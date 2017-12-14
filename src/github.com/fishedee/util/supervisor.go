package util

import (
	"errors"
	"github.com/fishedee/language"
	"os"
	"reflect"
	"sync"
)

var SupervisorStateEnum struct {
	language.EnumStruct
	STATE_STARTING int `enum:"1,开始中"`
	STATE_RUNNING  int `enum:"2,运行中"`
	STATE_STOPPING int `enum:"3,停止中"`
	STATE_EXITED   int `enum:"4,已停止"`
}

func init() {
	language.InitEnumStruct(&SupervisorStateEnum)
}

type SupervisorTaskInfo struct {
	stoppingChan chan bool
	exitChan     chan bool
	state        int
	taskId       int
}

type Supervisor struct {
	abortHandler func(taskId int)
	mutex        sync.Mutex
	task         map[int]*SupervisorTaskInfo
	startInner   func(taskInfo *SupervisorTaskInfo, name interface{}, argv ...interface{}) error
}

func newSupervisor() *Supervisor {
	supervisor := &Supervisor{}
	supervisor.task = map[int]*SupervisorTaskInfo{}
	return supervisor
}

func NewProcessSupervisor() *Supervisor {
	supervisor := newSupervisor()
	supervisor.startInner = func(taskInfo *SupervisorTaskInfo, name interface{}, argvs ...interface{}) error {
		nameString := name.(string)
		argvsString := []string{}
		for _, argv := range argvs {
			argvsString = append(argvsString, argv.(string))
		}
		process, err := os.StartProcess(nameString, argvsString, nil)
		if err != nil {
			return err
		}
		go func() {
			<-taskInfo.stoppingChan
			process.Kill()
		}()
		go func() {
			process.Wait()
			state := supervisor.GetState(taskInfo.taskId)
			if state == SupervisorStateEnum.STATE_RUNNING {
				go func() {
					supervisor.StopAndWait(taskInfo.taskId)
					if supervisor.abortHandler != nil {
						supervisor.abortHandler(taskInfo.taskId)
					}
				}()
			}
			taskInfo.exitChan <- true
		}()
		return nil
	}
	return supervisor
}

func NewThreadSupervisor() *Supervisor {
	supervisor := newSupervisor()
	supervisor.startInner = func(taskInfo *SupervisorTaskInfo, name interface{}, argvs ...interface{}) error {
		funValue := reflect.ValueOf(name)
		argvValue := []reflect.Value{}
		argvValue = append(argvValue, reflect.ValueOf(taskInfo.stoppingChan))
		for _, argv := range argvs {
			argvValue = append(argvValue, reflect.ValueOf(argv))
		}

		go func() {
			funValue.Call(argvValue)
			state := supervisor.GetState(taskInfo.taskId)
			if state == SupervisorStateEnum.STATE_RUNNING {
				go func() {
					//调用StopAndWait，是保证能清理taskId相关信息，以及在abortHandler回调时进程的状态是exited的
					supervisor.StopAndWait(taskInfo.taskId)
					if supervisor.abortHandler != nil {
						supervisor.abortHandler(taskInfo.taskId)
					}
				}()
			}
			taskInfo.exitChan <- true
		}()
		return nil
	}
	return supervisor
}

func (this *Supervisor) Start(taskId int, name interface{}, argv ...interface{}) error {
	this.mutex.Lock()
	taskInfo, isExist := this.task[taskId]
	if isExist == true {
		this.mutex.Unlock()
		return errors.New("已停止的任务才能启动!")
	}
	taskInfo = &SupervisorTaskInfo{}
	//stoppingChan设置为非阻塞，是为了实现threadSupervisor的异步stop
	taskInfo.stoppingChan = make(chan bool, 1)
	//exitChan设置为非阻塞，是为了让工作进程更快地退出，避免等待wait调用时才能退出
	taskInfo.exitChan = make(chan bool, 1)
	taskInfo.taskId = taskId
	taskInfo.state = SupervisorStateEnum.STATE_STARTING
	this.task[taskId] = taskInfo
	this.mutex.Unlock()

	err := this.startInner(taskInfo, name, argv...)
	if err != nil {
		this.mutex.Lock()
		delete(this.task, taskId)
		this.mutex.Unlock()
		return err
	}

	this.mutex.Lock()
	this.task[taskId].state = SupervisorStateEnum.STATE_RUNNING
	this.mutex.Unlock()

	return nil
}

func (this *Supervisor) wait(taskInfo *SupervisorTaskInfo) {
	_ = <-taskInfo.exitChan

	this.mutex.Lock()
	delete(this.task, taskInfo.taskId)
	this.mutex.Unlock()
}

func (this *Supervisor) stop(taskId int) (*SupervisorTaskInfo, error) {
	this.mutex.Lock()
	taskInfo, isExist := this.task[taskId]
	if isExist == false || taskInfo.state != SupervisorStateEnum.STATE_RUNNING {
		this.mutex.Unlock()
		return nil, errors.New("运行中的任务才能停止!")
	}
	this.task[taskId].state = SupervisorStateEnum.STATE_STOPPING
	this.mutex.Unlock()

	taskInfo.stoppingChan <- true
	return taskInfo, nil
}

func (this *Supervisor) Stop(taskId int) error {
	taskInfo, err := this.stop(taskId)
	if err != nil {
		return err
	}

	go this.wait(taskInfo)
	return nil
}

func (this *Supervisor) StopAndWait(taskId int) error {
	taskInfo, err := this.stop(taskId)
	if err != nil {
		return err
	}

	this.wait(taskInfo)
	return nil
}

func (this *Supervisor) StopAll() {
	allTaskId := []int{}
	this.mutex.Lock()
	for taskId, _ := range this.task {
		allTaskId = append(allTaskId, taskId)
	}
	this.mutex.Unlock()

	for _, taskId := range allTaskId {
		this.StopAndWait(taskId)
	}
}

func (this *Supervisor) GetState(taskId int) int {
	this.mutex.Lock()
	taskInfo, isExist := this.task[taskId]
	this.mutex.Unlock()

	if isExist == false {
		return SupervisorStateEnum.STATE_EXITED
	} else {
		return taskInfo.state
	}
}

func (this *Supervisor) GetBatchState(taskIds []int) map[int]int {
	result := map[int]int{}

	this.mutex.Lock()
	for _, taskId := range taskIds {
		var state int
		taskInfo, isExist := this.task[taskId]
		if isExist == false {
			state = SupervisorStateEnum.STATE_EXITED
		} else {
			state = taskInfo.state
		}
		result[taskId] = state
	}
	this.mutex.Unlock()

	return result
}

func (this *Supervisor) SetAbortHandler(abortHandler func(taskId int)) {
	this.abortHandler = abortHandler
}
