package workgroup

import (
	. "github.com/fishedee/app/log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type WorkGroupTask interface {
	Run() error
	Close()
}

type WorkGroup interface {
	Add(task WorkGroupTask)
	Run() error
	Close()
}

type WorkGroupConfig struct {
	CloseTimeout time.Duration `config:"closetimeout"`
	GraceClose   bool          `config:"graceclose"`
}

func NewWorkGroup(log Log, config WorkGroupConfig) (WorkGroup, error) {
	if config.CloseTimeout == 0 {
		config.CloseTimeout = time.Second * 5
	}
	return &workGroupImplement{
		log:       log,
		config:    config,
		closeChan: make(chan bool, 1),
	}, nil
}

type workGroupImplement struct {
	log       Log
	config    WorkGroupConfig
	task      []WorkGroupTask
	closeChan chan bool
}

func (this *workGroupImplement) Add(task WorkGroupTask) {
	this.task = append(this.task, task)
}

func (this *workGroupImplement) Run() error {
	errChan := make(chan error, len(this.task))
	doneChan := make(chan bool, 1)
	waitgroup := &sync.WaitGroup{}
	for i := 0; i != len(this.task); i++ {
		singleTask := this.task[i]
		waitgroup.Add(1)
		go func(singleTask WorkGroupTask) {
			defer waitgroup.Done()
			err := singleTask.Run()
			if err != nil {
				errChan <- err
			}
		}(singleTask)
	}

	go func() {
		waitgroup.Wait()
		doneChan <- true
	}()

	if this.config.GraceClose {
		this.setGraceClose()
	}
	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		this.log.Debug("workgroup is exited by self")
		return nil
	case <-this.closeChan:
		this.close(doneChan)
		this.log.Debug("workgroup is exited by close")
		return nil
	}
}

func (this *workGroupImplement) setGraceClose() {
	ch := make(chan os.Signal, 10)
	signals := []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGKILL,
	}
	signal.Notify(ch, signals...)
	go func() {
		<-ch
		this.Close()
	}()
}

func (this *workGroupImplement) close(doneChan chan bool) {
	closeFinishChan := make(chan bool)
	var lastClose int32
	go func() {
		for i := len(this.task) - 1; i >= 0; i-- {
			atomic.StoreInt32(&lastClose, int32(i))
			this.task[i].Close()
		}
		<-doneChan
		closeFinishChan <- true
	}()

	select {
	case <-time.After(this.config.CloseTimeout):
		this.log.Critical("workgroup wait %v because close but not exit, so force exit all,last not close work : %v", this.config.CloseTimeout, lastClose)
		return
	case <-closeFinishChan:
		return
	}
}

func (this *workGroupImplement) Close() {
	close(this.closeChan)
}
