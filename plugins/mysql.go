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
