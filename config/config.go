package config

type Config struct {
	// 应用标识
	Name string `yaml:"name"`
	// 应用模式 dev/pro
	Mode string `yaml:"mode"`
	// 应用版本
	Version string `yaml:"version"`
	// 监听地址
	Listen string `yaml:"listen"`
	// 日志级别
	Level string `yaml:"level"`
	// 日志配置
	Logger string `yaml:"logger"`
	// MYSQL配置
	Mysql string `yaml:"mysql"`
	// SQLITE配置
	Sqlite string `yaml:"sqlite"`
}

func New() *Config {
	return &Config{
		Name: "app",
		Mode: "pro",
		Version: "0.0.0",
		Listen: ":8800",
		Level: "debug",
		Logger: "./config/logger.yaml",
		Mysql: "./config/mysql.yaml",
		Sqlite: "./config/sqlite.yaml",
	}
}
