package application

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	hooks := LoadFromEnv(AppConfigPrefix)

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
	logger := app.logger.
		Named("handler").
		With(
			zap.String("ip", req.RemoteAddr),
			zap.String("host", req.Host),
			zap.String("path", req.URL.Path),
		)
	logger.Info("Beginning handling")

	hooks := app.hooks.FindAllMatching(req.URL.Path)
	if len(hooks) == 0 {
		logger.Info("Could not match request")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not match request"))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(hooks))
	hookNames := make([]string, 0, len(hooks))
	for _, hook := range hooks {
		hookNames = append(hookNames, hook.Name)
		go app.proxy(req, hook, &wg)
	}
	wg.Wait()

	w.Header().Set("X-Hooks-Count", strconv.Itoa(len(hookNames)))
	w.Header().Set("X-Hooks", strings.Join(hookNames, "; "))
	w.WriteHeader(http.StatusCreated)

	// We don't wait for the results of all hooks execution - we don't care
	logger.Info("Finished handling")
}
