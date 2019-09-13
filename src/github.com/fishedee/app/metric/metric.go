package metric

import (
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

type MetricConfig struct {
	ConnectUrl string `config:"driver"`
	Database   string `config:"database"`
	User       string `config:"user"`
	Password   string `config:"password"`
}
type Metric interface {
	SetDefaultTags(defaultTags map[string]string)
	IncCounter(name string, counter int64)
	UpdateGauge(name string, data int64)
	UpdateHistogram(name string, data int64)
	MarkMeter(name string, data int64)
	UpdateTimer(name string, duration time.Duration)

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

func (this *metricImplement) IncCounter(name string, data int64) {
	counter := metrics.GetOrRegisterCounter(name, this.registry)

	counter.Inc(data)
}

func (this *metricImplement) UpdateGauge(name string, data int64) {
	gauge := metrics.GetOrRegisterGauge(name, this.registry)

	gauge.Update(data)
}

func (this *metricImplement) UpdateGaugeFloat64(name string, data float64) {
	gauge := metrics.GetOrRegisterGaugeFloat64(name, this.registry)

	gauge.Update(data)
}

func (this *metricImplement) UpdateHistogram(name string, data int64) {
	s := metrics.NewExpDecaySample(1028, 0.015)

	histogram := metrics.GetOrRegisterHistogram(name, this.registry, s)

	histogram.Update(data)
}

func (this *metricImplement) MarkMeter(name string, data int64) {
	meter := metrics.GetOrRegisterMeter(name, this.registry)

	meter.Mark(data)
}

func (this *metricImplement) UpdateTimer(name string, data time.Duration) {
	timer := metrics.GetOrRegisterTimer(name, this.registry)

	timer.Update(data)
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
