package application

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"http-hook-gateway/internal/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	AppConfigPrefix = "hhgw"
)

var ErrInitError = errors.New("init error")

type app struct {
	hooks      HookMap
	httpClient *http.Client
	logger     *zap.Logger
}

func NewApp(logger *zap.Logger) *app {
	return &app{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (app *app) LoadConfig() error {
	hooks := config.LoadFromEnv(AppConfigPrefix)

	if len(hooks) == 0 {
		app.logger.Fatal("No hooks loaded, nowhere to proxy")
		return ErrInitError
	}

	app.logger.Info("Loaded hooks", zap.Int("count", len(hooks)), zap.Strings("names", hooks.PeekNames()))

	parseErrs := hooks.ParseHooks()
	if len(parseErrs) != 0 {
		app.logger.Error("error parsing hooks", zap.Errors("errors", parseErrs))
		return ErrInitError
	}

	app.hooks = hooks

	return nil
}

func (app *app) RootHandler(w http.ResponseWriter, req *http.Request) {
	if app.hooks == nil {
		app.logger.Error("app hooks are not initialized, unable to continue")
		return
	}

	hooks := app.hooks.FindAllMatching(req.URL.Path)
	if len(hooks) == 0 {
		app.logger.Info("Could not match request", zap.String("path", req.URL.Path))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not match request"))
		return
	}

	for _, hook := range hooks {
		go app.proxy(req, hook)
	}

	// We don't wait for the results of all hooks execution - we don't care
	app.logger.Info("HTTP request handler exited")
}

func (app *app) proxy(req *http.Request, hookDef *HookDef) (Err error) {
	logger := app.logger.Named(fmt.Sprintf("%v %v", hookDef.Name, req.URL.Path))
	logger = logger.With(zap.Any("hook", hookDef))

	defer func() {
		if Err != nil {
			logger.Error("error", zap.Error(Err))
		} else {
			logger.Info("end")
		}
	}()
	logger.Info("begin")

	hook, parseErr := hookDef.ParseRequest(req)

	if parseErr != nil {
		return parseErr
	}

	url := fmt.Sprintf("%v/%v",
		strings.TrimRight(hookDef.ProxyHost, "/"),
		strings.TrimRight(hook.Path(hookDef.ProxyPath), "/"))

	logger.Info(
		"parsed request",
		zap.String("method", hook.Method()),
		zap.String("url", url),
		zap.Any("headers", hook.Headers()),
	)

	body, _ := ioutil.ReadAll(hook.Body())
	logger.Debug("body", zap.ByteString("body", body))

	proxyReq, reqCreateErr := http.NewRequest(hook.Method(), url, bytes.NewBuffer(body))
	if reqCreateErr != nil {
		return reqCreateErr
	}

	for name, val := range hook.Headers() {
		proxyReq.Header.Set(name, val)
	}

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
