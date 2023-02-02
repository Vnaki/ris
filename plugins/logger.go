package plugins

import (
	"fmt"
	"github.com/vnaki/ris/components/logger"
	"github.com/vnaki/ris/types"
)

func LoggerPlugin(name string, e types.Engine) error {
	l := logger.New()

	if err := e.Parse(e.Config().Logger, l); err != nil {
		return err
	}

	log := e.App().Logger()

	if !e.IsDev() {
		fd, err := l.Open()
		if err != nil {
			return err
		}

		log.Handle(l.Handle)
		log.SetOutput(fd)
	}

	e.Set(name, log)

	e.Defer(func() {
		l.Close()

		// verbose
		fmt.Println("defer: log closed")
	})

	return nil
}
