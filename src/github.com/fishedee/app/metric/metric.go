package metric

import (
	"github.com/rcrowley/go-metrics"
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
	SetDefaultTags(defaultTags map[string]string)
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
	registry    metrics.Registry
	config      MetricConfig
	defaultTags map[string]string
	closeChan   chan bool
	waitgroup   *sync.WaitGroup
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

func (this *metricImplement) SetDefaultTags(defaultTags map[string]string) {
	this.defaultTags = defaultTags
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
			this.defaultTags,
			this.closeChan,
		)
	}()

	this.waitgroup.Wait()

	return nil
}

func (this *metricImplement) Close() {
	close(this.closeChan)
}
