package config

import "github.com/kataras/iris/v12"

// IrisConfig iris配置
var IrisConfig = iris.Configuration{
	// 是否支持重复绑定端口
	SocketSharding: true,
	// 是否禁用启动日志
	DisableStartupLog: false,
	// 是否禁用中断(CTRL+C)
	DisableInterruptHandler: false,
	// 是否禁用路径矫正(尾部斜杠), /home/ -> /home
	DisablePathCorrection: false,
	// 是否禁用重定向路径矫正
	DisablePathCorrectionRedirection: true,
	// 是否开辟新的缓冲区接受数据而不是从`context.UnmarshalBody/ReadJSON/ReadXML`消费数据
	DisableBodyConsumptionOnUnmarshal: false,
	// 是否重定向不存在到接近的路由, /con -> /constants
	EnablePathIntelligence: false,
	// 是否路径和参数自动转义
	EnablePathEscape: true,
	// 尝试优化程序
	EnableOptimizations: true,
	// 是否将不支持的方法以405代替404
	FireMethodNotAllowed: false,
	// 是否报错`iris.ErrEmptyForm`,如果读`context.ReadBody/ReadForm`空数据
	FireEmptyFormError: false,
	// 最大POST数据尺寸, 20MB
	PostMaxMemory: 20 << 20,
}
