package ris

import (
	"context"
	"errors"
	"fmt"
	"github.com/vnaki/ris/config"
	"github.com/vnaki/ris/constants"
	"github.com/vnaki/ris/types"
	"net/http"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"gopkg.in/yaml.v3"
)

type Engine struct {
	app         *iris.Application
	conf        *config.Config
	reset       func(*config.Config)
	modules     map[string]types.Module
	workers     map[string]types.Worker
	components  map[string]types.Component
	middlewares map[string]types.Middleware
	plugins     []types.Plugin
	implements  []func() error
	deferments  []func()
}

func New() *Engine {
	return &Engine{
		app:         iris.New(),
		conf:        config.New(),
		modules:     make(map[string]types.Module),
		workers:     make(map[string]types.Worker),
		components:  make(map[string]types.Component),
		middlewares: make(map[string]types.Middleware),
	}
}

func (e *Engine) App() *iris.Application {
	return e.app
}

func (e *Engine) Config() *config.Config {
	return e.conf
}

func (e *Engine) Set(name string, component types.Component) {
	e.components[name] = component
}

func (e *Engine) SetPostMemory(memory int64) {
	config.IrisConfig.PostMaxMemory = memory
}

func (e *Engine) Get(name string) types.Component {
	if value, ok := e.components[name]; ok {
		return value
	}
	return nil
}

func (e *Engine) Reset(f func(c *config.Config)) {
	if e.reset == nil {
		e.reset = f
	}
}

func (e *Engine) Plugin(name string, plugin types.PluginHandler) {
	e.plugins = append(e.plugins, types.Plugin{
		Name:    name,
		Handler: plugin,
	})
}

func (e *Engine) Worker(name string, worker types.Worker) {
	e.workers[name] = worker
}

func (e *Engine) Module(name string, module types.Module) {
	e.modules[name] = module
}

func (e *Engine) Middleware(middleware types.MiddlewareHandler) {
	e.IfMiddleware("*", middleware)
}

func (e *Engine) IfMiddleware(env string, middleware types.MiddlewareHandler) {
	e.middlewares[env] = types.Middleware{
		Route:   false,
		Handler: middleware,
	}
}

func (e *Engine) RouteMiddleware(middleware types.MiddlewareHandler) {
	e.IfRouteMiddleware("*", middleware)
}

func (e *Engine) IfRouteMiddleware(env string, middleware types.MiddlewareHandler) {
	e.middlewares[env] = types.Middleware{
		Route:   true,
		Handler: middleware,
	}
}

func (e *Engine) Stop() error {
	return e.app.Shutdown(context.Background())
}

func (e *Engine) Implement(fn func() error) {
	e.implements = append(e.implements, fn)
}

func (e *Engine) Parse(file string, out interface{}) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file)
	}

	b, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read file error: %s, file: %s", err.Error(), file)
	}

	if err = yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("unmarshal file error: %s, file: %s", err.Error(), file)
	}

	return nil
}

func (e *Engine) IsDev() bool {
	return e.conf.Mode != "pro"
}

func (e *Engine) Defer(f func()) {
	e.deferments = append(e.deferments, f)
}

func (e *Engine) Test(file string) error {
	if err := e.Parse(file, e.conf); err != nil {
		return err
	}

	if e.reset != nil {
		e.reset(e.conf)
	}

	for _, plugin := range e.plugins {
		e.app.Logger().Debugf("register `%s` plugin", plugin.Name)

		if err := plugin.Handler(plugin.Name, e); err != nil {
			return fmt.Errorf("register `%s` plugin failed: %v", plugin.Name, err.Error())
		}
	}

	for _, implement := range e.implements {
		if err := implement(); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Run(file string) error {
	if err := e.Parse(file, e.conf); err != nil {
		return err
	}

	if e.reset != nil {
		e.reset(e.conf)
	}

	e.app.SetName(e.conf.Name)
	e.app.Logger().SetLevel(e.conf.Level)
	e.app.Logger().SetTimeFormat(constants.TimeFormat)
	e.app.Use(recover.New())

	for mode, middleware := range e.middlewares {
		if mode == "*" || e.conf.Mode == mode {
			if middleware.Route {
				e.app.UseRouter(middleware.Handler(e))
			} else {
				e.app.Use(middleware.Handler(e))
			}
		}
	}

	for _, plugin := range e.plugins {
		e.app.Logger().Debugf("register `%s` plugin", plugin.Name)

		if err := plugin.Handler(plugin.Name, e); err != nil {
			return fmt.Errorf("register `%s` plugin failed: %v", plugin.Name, err.Error())
		}
	}

	for name, worker := range e.workers {
		go worker(name, e)
	}

	for name, module := range e.modules {
		mvc.Configure(e.app.Party(name), module)
	}

	e.app.OnBuild = func() error {
		for _, implement := range e.implements {
			if err := implement(); err != nil {
				return err
			}
		}
		return nil
	}

	return e.serve()
}

func (e *Engine) serve() error {
	defer func() {
		for _, f := range e.deferments {
			f()
		}
	}()

	if !e.IsDev() {
		config.IrisConfig.DisableStartupLog = true
	}

	if err := e.app.Run(
		iris.Addr(e.conf.Listen),
		iris.WithConfiguration(config.IrisConfig),
	); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("app(%v) boot failed: %v", e.conf.Name, err.Error())
	}

	return nil
}
