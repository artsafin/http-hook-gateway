package application

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	UserAgent = "hhgw/1.0"
)

func (app *app) proxy(req *http.Request, hookDef *HookDef, wg *sync.WaitGroup) (Err error) {
	logger := app.logger.Named(fmt.Sprintf("%v %v", hookDef.Name, req.URL.Path))

	defer func() {
		if Err != nil {
			logger.
				WithOptions(zap.AddCallerSkip(1)).
				Error("error", zap.Error(Err))
		}

		wg.Done()
	}()
	logger.Info("begin", zap.String("proxy_host", hookDef.ProxyHost))

	hook, parseErr := hookDef.ParseRequest(req)

	if parseErr != nil {
		return parseErr
	}

	url := fmt.Sprintf("%v/%v",
		strings.TrimRight(hookDef.ProxyHost, "/"),
		strings.TrimLeft(hook.Path(hookDef.ProxyPath), "/"))

	logger.Info(
		"parsed request",
		zap.String("method", hook.Method()),
		zap.String("url", url),
		zap.Any("headers", hook.Headers()),
	)

	body, _ := ioutil.ReadAll(hook.Body())

	proxyReq, reqCreateErr := http.NewRequest(hook.Method(), url, bytes.NewBuffer(body))
	if reqCreateErr != nil {
		return reqCreateErr
	}

	for name, val := range hook.Headers() {
		proxyReq.Header.Set(name, val)
	}
	proxyReq.Header.Set("User-Agent", UserAgent)

	resp, reqErr := app.httpClient.Do(proxyReq)
	if reqErr != nil {
		return reqErr
	}

	respBodyBytes, respErr := ioutil.ReadAll(resp.Body)
	// Non fatal
	if respErr != nil {
		logger.Error("error reading response body", zap.Error(respErr))
	}
	defer resp.Body.Close()

	var respHeaders strings.Builder
	_ = resp.Header.Write(&respHeaders)

	logger.Info(
		"response",
		zap.String("status", resp.Status),
		zap.String("headers", respHeaders.String()),
		zap.ByteString("body", respBodyBytes),
	)

	return nil
}
