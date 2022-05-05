//修改自https://github.com/influxdata/influxdb1-client
//原来代码对schema的设计很不合理
package metric

import (
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/rcrowley/go-metrics"
	"log"
	"strings"
	"time"
	"math"
)

type reporter struct {
	reg      metrics.Registry
	interval time.Duration
	url      string
	database string
	username string
	password string
	client   client.Client
	exitChan chan bool
}

func InfluxDB(r metrics.Registry, d time.Duration, url, database, username, password string, exitChan chan bool) {
	rep := &reporter{
		reg:      r,
		interval: d,
		url:      url,
		database: database,
		username: username,
		password: password,
		exitChan: exitChan,
	}
	if err := rep.makeClient(); err != nil {
		panic(fmt.Sprintf("unable to make InfluxDB client. err=%v", err))
	}

	rep.run()
}

func (r *reporter) makeClient() error {
	var err error
	r.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     r.url,
		Username: r.username,
		Password: r.password,
	})

	return err
}

func (r *reporter) run() {
	intervalTicker := time.Tick(r.interval)
	pingTicker := time.Tick(time.Second * 5)
	for {
		select {
		case <-intervalTicker:
			if err := r.send(); err != nil {
				log.Printf("unable to send metrics to InfluxDB. err=%v", err)
			}
		case <-pingTicker:
			_, _, err := r.client.Ping(time.Second)
			if err != nil {
				log.Printf("got error while sending a ping to InfluxDB, trying to recreate client. err=%v", err)

				if err = r.makeClient(); err != nil {
					log.Printf("unable to make InfluxDB client. err=%v", err)
				}
			}
		case <-r.exitChan:
			return
		}
	}
}

func (r *reporter) getTags(name string) (string, map[string]string) {
	result := map[string]string{}
	leftBraceIndex := strings.IndexByte(name, '?')
	if leftBraceIndex == -1 {
		return name, result
	}
	measurement := name[0:leftBraceIndex]
	tagStr := name[leftBraceIndex+1:]
	tagList := Explode(tagStr, "&")
	for _, tag := range tagList {
		tagInfo := Explode(tag, "=")
		if len(tagInfo) == 2 {
			var err error
			result[tagInfo[0]], err = DecodeUrl(tagInfo[1])
			if err != nil {
				panic(err)
			}
		}
	}
	return measurement, result
}

func (r *reporter) getRate( data float64) float64{
	if math.IsInf(data,0) || math.IsNaN(data){
		return 0
	}else{
		return data
	}
}

func (r *reporter) send() error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: r.database,
	})
	if err != nil {
		return err
	}

	now := time.Now().Truncate(r.interval)

	r.reg.Each(func(name string, i interface{}) {
		measurement, tags := r.getTags(name)

		switch metric := i.(type) {
		case metrics.Counter:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				"count": ms.Count(),
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)
		case metrics.Gauge:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				"gauge": ms.Value(),
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)
		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				"gauge": ms.Value(),
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)
		case metrics.Histogram:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			fields := map[string]interface{}{
				"count":    ms.Count(),
				"max":      ms.Max(),
				"mean":     ms.Mean(),
				"min":      ms.Min(),
				"stddev":   ms.StdDev(),
				"variance": ms.Variance(),
				"p50":      ps[0],
				"p75":      ps[1],
				"p95":      ps[2],
				"p99":      ps[3],
				"p999":     ps[4],
				"p9999":    ps[5],
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)
		case metrics.Meter:
			ms := metric.Snapshot()
			fields := map[string]interface{}{
				"count":    ms.Count(),
				"m1":       ms.Rate1(),
				"m5":       ms.Rate5(),
				"m15":      ms.Rate15(),
				"meanrate": r.getRate(ms.RateMean()),
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)
		case metrics.Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			fields := map[string]interface{}{
				"count":    ms.Count(),
				"max":      ms.Max(),
				"mean":     ms.Mean(),
				"min":      ms.Min(),
				"stddev":   ms.StdDev(),
				"variance": ms.Variance(),
				"p50":      ps[0],
				"p75":      ps[1],
				"p95":      ps[2],
				"p99":      ps[3],
				"p999":     ps[4],
				"p9999":    ps[5],
				"m1":       ms.Rate1(),
				"m5":       ms.Rate5(),
				"m15":      ms.Rate15(),
				"meanrate": r.getRate(ms.RateMean()),
			}
			pt, err := client.NewPoint(measurement, tags, fields, now)
			if err != nil {
				panic(err)
			}
			bp.AddPoint(pt)

		default:
			panic(fmt.Sprintf("unknown metric type %v", metric))
		}
	})
	err = r.client.Write(bp)
	return err
}
