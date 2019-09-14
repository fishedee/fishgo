package metric

import (
	"fmt"
	"testing"
	"time"
)

func TestMetric(t *testing.T) {
	metric, err := NewMetric(MetricConfig{
		ConnectUrl: "http://localhost:8086",
		Database:   "test",
		User:       "",
		Password:   "",
	})
	if err != nil {
		panic(err)
	}

	go func() {
		counter := metric.GetCounter("reqTime")
		timer := metric.GetTimer("reqTime2")
		gauge := metric.GetGauge("reqTime3?path=/user/get")

		for i := 0; i != 10; i++ {
			time.Sleep(time.Millisecond * 100)
			//递增计数器
			counter.Inc(1)

			//计量时间与次数
			timer.Update(time.Second * 2)

			//计量Gauge
			gauge.Update(12)

			fmt.Println("finish")
		}
		metric.Close()
	}()

	metric.Run()
}
