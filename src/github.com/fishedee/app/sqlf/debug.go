package sqlf

import (
	. "github.com/fishedee/app/log"
	"time"
)

func runSql(isDebug bool, log Log, handler func() (string, error)) error {
	if isDebug {
		beginTime := time.Now()
		sql, err := handler()
		duration := time.Now().Sub(beginTime)
		log.Debug("[sqlf] sql:[%s] isErr:[%v] duration:[%v]", sql, err, duration)
		return err
	} else {
		_, err := handler()
		return err
	}
}
