package logger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/kataras/golog"
)

type Logger struct {
	// 链接名称
	Link string `yaml:"link"`
	// 储存路径
	Path string `yaml:"path"`
	// 文件权限
	Perm fs.FileMode `yaml:"perm"`
	// 是否按天切割日志
	Daily bool `yaml:"daily"`
	// 单个文件大小
	Size int64 `yaml:"size"`
	// 文件名前缀
	Prefix string `yaml:"prefix"`
	// 日期模板
	layout string `yaml:"layout"`

	fd *os.File
}

func New() *Logger {
	return &Logger{
		Path:   "./logs",
		Link:   "latest",
		Daily:  true,
		Size:   1024,
		Perm:   0666,
		Prefix: "log",
		layout: "20060102",
	}
}

func (l *Logger) Open() (*os.File, error) {
	if err := os.MkdirAll(l.Path, os.ModePerm); err != nil {
		return nil, err
	}

	file := fmt.Sprintf("%s/%s.log", l.Path, l.Link)

	fd, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, l.Perm)
	if err != nil {
		return nil, err
	}

	l.fd = fd

	return l.fd, nil
}

func (l *Logger) Close() {
	if l.fd != nil {
		_ = l.fd.Close()
	}
}

func (l *Logger) Handle(*golog.Log) bool {
	info, err := l.fd.Stat()
	if err != nil {
		panic(fmt.Sprintf("logger error, file:%s, err:%s", l.fd.Name(), err.Error()))
	}

	file, err := l.path(info.ModTime())
	if err != nil {
		panic(err)
	}

	if (l.Daily && time.Now().Format(l.layout) != info.ModTime().Format(l.layout)) || info.Size() >= l.Size {
		if err := l.copy(file); err != nil {
			panic("logger copy error: " + err.Error())
		}
	}

	return false
}

func (l *Logger) copy(file string) error {
	dst, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, l.Perm)
	if err != nil {
		return err
	}

	defer func() {
		_ = dst.Close()
	}()

	if _, err = l.fd.Seek(0, 0); err != nil {
		return err
	}

	if _, err = io.Copy(dst, l.fd); err != nil && err != io.EOF {
		return err
	}

	return l.fd.Truncate(0)
}

func (l *Logger) path(t time.Time) (string, error) {
	path := l.Path

	if l.Prefix != "" {
		path += fmt.Sprintf("/%s-", l.Prefix)
	}

	if l.Daily {
		path += t.Format(l.layout) + "-"
	}

	if matches, err := filepath.Glob(path + "*.log"); err != nil {
		return "", err
	} else {
		path += fmt.Sprintf("%02d.log", len(matches)+1)
	}

	return path, nil
}
