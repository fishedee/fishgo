package sqlf

import (
	"fmt"
	"time"
)

func runSql(isDebug bool, handler func() (string, error)) error {
	if isDebug {
		beginTime := time.Now()
		sql, err := handler()
		duration := time.Now().Sub(beginTime)
		fmt.Printf("[sqlf] sql:%s isErr:%v,duration:%v", sql, err, duration)
		return err
	} else {
		_, err := handler()
		return err
	}
}
