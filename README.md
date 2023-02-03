# ris - 基于iris的插件式框架

**Iris framework for plug-in design, super easy to use！**

### 快速开始

```go 

package main

import (
	"github.com/vnaki/ris"
	"github.com/vnaki/ris/examples/routes"
	"github.com/vnaki/ris/middlewares"
	"github.com/vnaki/ris/plugins"
)

func main()  {
	e := ris.New()

	// post max memory
	e.SetPostMemory(20 << 20)

	e.RouteMiddleware(middlewares.Cors)

	e.Plugin("logger", plugins.LoggerPlugin)
	//e.Plugin("data", plugins.MysqlPlugin)
	e.Plugin("data", plugins.SqlitePlugin)

	// default module
	e.Module("/", routes.ApiRoute)

	if err := e.Run("./config/app.yaml"); err != nil {
		panic(err)
	}
}

```

### 详细用例

see [examples](https://github.com/vnaki/ris/tree/master/examples)

### 框架定义

```go 
package types

import (
	"github.com/vnaki/ris/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Database 数据库接口
type Database interface {
	Connect() (db.Session, error)
}

// PluginHandler 插件函数
type PluginHandler func(string, Engine) error

// Plugin 插件结构
type Plugin struct {
	// 插件名称
	Name string
	// 插件函数
	Handler PluginHandler
}

// Module 模块函数
type Module func(*mvc.Application)

// MiddlewareHandler 中间件处理
type MiddlewareHandler func(Engine) iris.Handler

// Middleware 中间件
type Middleware struct {
	// 是否路由中间件
	Route bool
	// 中间件处理
	Handler MiddlewareHandler
}

// Component 组件类型
type Component interface{}

// Worker 工作协程
type Worker func(string, Engine)

// Engine 应用引擎
type Engine interface {
	// App iris应用
	App() *iris.Application
	// Set 注册实例
	Set(name string, component Component)
    // SetPostMemory 设置POST最大内存
	SetPostMemory(memory int64)
	// Get 返回实例
	Get(name string) Component
	// Reset 重置配置
	Reset(func(c *config.Config))
	// Config 配置信息
	Config() *config.Config
	// Plugin 注册插件
	Plugin(name string, plugin PluginHandler)
	// Worker 注册工作协程
	Worker(name string, plugin Worker)
	// Module 注册模块
	Module(party string, module Module)
	// Middleware 注册通用中间件
	Middleware(handler MiddlewareHandler)
	// RouteMiddleware 路由中间件
	RouteMiddleware(middleware MiddlewareHandler)
	// IfMiddleware 注册条件中间件
	IfMiddleware(mode string, middleware MiddlewareHandler)
	// IfRouteMiddleware 注册路由条件中间件
	IfRouteMiddleware(mode string, middleware MiddlewareHandler)
	// Stop 停止服务
	Stop() error
	// Implement 自定义业务
	Implement(func() error)
	// Defer 延迟函数
	Defer(f func())
	// IsDev 是否开发环境
	IsDev() bool
	// Run 运行程序
	Run(file string) error
	// Test 执行测试
	Test(file string) error
	// Parse 解析配置
	Parse(file string, out interface{}) error
}

```

### 插件定义

```go 
package plugins

import (
	"fmt"
	"github.com/vnaki/ris/components/database"
	"github.com/vnaki/ris/types"
)

func MysqlPlugin(name string, e types.Engine) error {
	n := database.New()

	if err := e.Parse(e.Config().Mysql, n); err != nil {
		return err
	}

	sess, err := n.Connect()
	if err != nil {
		return err
	}

	e.Set(name, sess)

	e.Defer(func() {
		_ = sess.Close()

		// verbose
		fmt.Println("defer: mysql closed")
	})

	return nil
}

```

### 框架依赖

- `iris` framework, see [https://github.com/kataras/iris](https://github.com/kataras/iris)
- `upper` orm, see [https://github.com/upper/db](https://github.com/upper/db)
