package workgroup

import (
	. "github.com/fishedee/app/log"
	"sync"
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
}

func NewWorkGroup(log Log, config WorkGroupConfig) (WorkGroup, error) {
	if config.CloseTimeout == 0 {
		config.CloseTimeout = time.Second * 5
	}
	return &workGroupImplement{
		log:       log,
		config:    config,
		hasClose:  false,
		closeChan: make(chan bool, 1),
	}, nil
}

type workGroupImplement struct {
	log       Log
	config    WorkGroupConfig
	task      []WorkGroupTask
	loadTask  sync.Map
	mutex     sync.Mutex
	hasClose  bool
	closeChan chan bool
}

func (this *workGroupImplement) Add(task WorkGroupTask) {
	this.task = append(this.task, task)
}

func (this *workGroupImplement) waitDoneOrTimeout(doneChan chan bool) {
	select {
	case <-time.After(this.config.CloseTimeout):
		this.log.Critical("workgroup wait %v because close but not exit, so force exit all", this.config.CloseTimeout)
		return
	case <-doneChan:
		return
	}
}
func (this *workGroupImplement) Run() error {
	errChan := make(chan error, len(this.task))
	doneChan := make(chan bool, 1)
	waitgroup := &sync.WaitGroup{}
	for i := 0; i != len(this.task); i++ {
		singleTask := this.task[i]
		this.loadTask.Store(singleTask, true)
		waitgroup.Add(1)
		go func(singleTask WorkGroupTask) {
			err := singleTask.Run()
			if err != nil {
				errChan <- err
			}
			this.loadTask.Delete(singleTask)
			waitgroup.Done()
		}(singleTask)
	}

	go func() {
		waitgroup.Wait()
		doneChan <- true
	}()
	select {
	case err := <-errChan:
		this.Close()
		this.waitDoneOrTimeout(doneChan)
		return err
	case <-doneChan:
		return nil
	case <-this.closeChan:
		this.waitDoneOrTimeout(doneChan)
		return nil
	}
}

func (this *workGroupImplement) Close() {
	this.mutex.Lock()
	if this.hasClose == true {
		this.mutex.Unlock()
		return
	}
	this.closeChan <- true
	this.hasClose = true
	this.mutex.Unlock()

	this.loadTask.Range(func(key interface{}, value interface{}) bool {
		singleTask := key.(WorkGroupTask)
		go func() {
			singleTask.Close()
		}()
		return true
	})
}
