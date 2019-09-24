package metric

import (
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"github.com/rcrowley/go-metrics"
	"strings"
	"sync"
	"time"
)

type MetricCounter = metrics.Counter

type MetricGauge = metrics.Gauge

type MetricGaugeFloat64 = metrics.GaugeFloat64

type MetricHistogram = metrics.Histogram

type MetricMeter = metrics.Meter

type MetricTimer = metrics.Timer

type MetricConfig struct {
	ConnectUrl string `config:"connecturl"`
	Database   string `config:"database"`
	User       string `config:"user"`
	Password   string `config:"password"`
}

type Metric interface {
	GetCounter(name string) MetricCounter
	GetGauge(name string) MetricGauge
	GetGaugeFloat64(name string) MetricGaugeFloat64
	GetHistogram(name string) MetricHistogram
	GetMeter(name string) MetricMeter
	GetTimer(name string) MetricTimer

	Run() error
	Close()
}

type metricImplement struct {
	registry  metrics.Registry
	config    MetricConfig
	closeChan chan bool
	waitgroup *sync.WaitGroup
}

func NewMetric(config MetricConfig) (Metric, error) {
	registry := metrics.NewRegistry()
	return &metricImplement{
		config:    config,
		registry:  registry,
		closeChan: make(chan bool),
		waitgroup: &sync.WaitGroup{},
	}, nil
}

func (this *metricImplement) GetCounter(name string) metrics.Counter {
	return metrics.GetOrRegisterCounter(name, this.registry)
}

func (this *metricImplement) GetGauge(name string) metrics.Gauge {
	return metrics.GetOrRegisterGauge(name, this.registry)
}

func (this *metricImplement) GetGaugeFloat64(name string) metrics.GaugeFloat64 {
	return metrics.GetOrRegisterGaugeFloat64(name, this.registry)
}

func (this *metricImplement) GetHistogram(name string) metrics.Histogram {
	s := metrics.NewExpDecaySample(1028, 0.015)
	return metrics.GetOrRegisterHistogram(name, this.registry, s)
}

func (this *metricImplement) GetMeter(name string) metrics.Meter {
	return metrics.GetOrRegisterMeter(name, this.registry)
}

func (this *metricImplement) GetTimer(name string) metrics.Timer {
	return metrics.GetOrRegisterTimer(name, this.registry)
}

func (this *metricImplement) Run() error {
	metrics.RegisterDebugGCStats(this.registry)
	metrics.RegisterRuntimeMemStats(this.registry)

	this.waitgroup.Add(2)

	//启动runtime的debug上传
	go func() {
		defer this.waitgroup.Done()

		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				metrics.CaptureDebugGCStatsOnce(this.registry)
				metrics.CaptureRuntimeMemStatsOnce(this.registry)
				break
			case <-this.closeChan:
				return
			}
		}
	}()

	//启动埋点提交influxDb
	go func() {
		defer this.waitgroup.Done()

		InfluxDB(this.registry,
			time.Second,
			this.config.ConnectUrl,
			this.config.Database,
			this.config.User,
			this.config.Password,
			this.closeChan,
		)
	}()

	this.waitgroup.Wait()

	return nil
}

func (this *metricImplement) Close() {
	close(this.closeChan)
}

type metricTagsImplement struct {
	metric     Metric
	taggedName string
}

func (this *metricTagsImplement) GetCounter(name string) MetricCounter {
	return this.metric.GetCounter(this.getName(name))
}

func (this *metricTagsImplement) GetGauge(name string) MetricGauge {
	return this.metric.GetGauge(this.getName(name))
}

func (this *metricTagsImplement) GetGaugeFloat64(name string) MetricGaugeFloat64 {
	return this.metric.GetGaugeFloat64(this.getName(name))
}

func (this *metricTagsImplement) GetHistogram(name string) MetricHistogram {
	return this.metric.GetHistogram(this.getName(name))
}

func (this *metricTagsImplement) GetMeter(name string) MetricMeter {
	return this.metric.GetMeter(this.getName(name))
}

func (this *metricTagsImplement) GetTimer(name string) MetricTimer {
	return this.metric.GetTimer(this.getName(name))
}

func (this *metricTagsImplement) Run() error {
	return this.metric.Run()
}

func (this *metricTagsImplement) Close() {
	this.metric.Close()
}

func (this *metricTagsImplement) getName(data string) string {
	questionIndex := strings.IndexByte(data, '?')
	if questionIndex == -1 {
		return data + "?" + this.taggedName
	} else {
		return data + "&" + this.taggedName
	}
}

func getTaggedName(tags map[string]string) string {
	tagList := []string{}
	if tags != nil {
		for k, v := range tags {
			vEncode, err := EncodeUrl(v)
			if err != nil {
				panic(err)
			}
			tagList = append(tagList, k+"="+vEncode)
		}
	}
	tagStr := Implode(tagList, "&")
	return tagStr
}

func MetricWithDefaultTags(metric Metric, tags map[string]string) Metric {
	taggedName := getTaggedName(tags)
	return &metricTagsImplement{
		metric:     metric,
		taggedName: taggedName,
	}
}
