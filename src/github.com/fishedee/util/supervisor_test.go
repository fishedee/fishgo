package util

import (
	. "github.com/fishedee/assert"
	"runtime"
	"testing"
	"time"
)

//检查routine数量
func AssertRoutineOne(t *testing.T) {
	if runtime.NumGoroutine() != 2 {
		panic("AssertRoutineOne Fail!")
	}
}

//测试输入参数
func TestSupervisorArgv(t *testing.T) {
	supervisor := NewThreadSupervisor()

	var argvInt int
	handler := func(stopChan chan bool, input int) {
		argvInt = input
		<-stopChan
	}

	//测试输入参数
	supervisor.Start(1, handler, 100)
	supervisor.StopAndWait(1)
	AssertEqual(t, argvInt, 100, "input argv != 100")
	AssertRoutineOne(t)
}

//测试多个实例启动
func TestSupervisorMultiStart(t *testing.T) {
	supervisor := NewThreadSupervisor()

	var argvInt int
	handler := func(stopChan chan bool, input int) {
		<-stopChan
		argvInt = input
	}
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state1")
	supervisor.Start(1, handler, 100)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state2")

	var argvInt2 int
	handler2 := func(stopChan chan bool, input int) {
		<-stopChan
		argvInt2 = input
	}

	AssertEqual(t, supervisor.GetState(2), SupervisorStateEnum.STATE_EXITED, "state3")
	supervisor.Start(2, handler2, 101)
	AssertEqual(t, supervisor.GetState(2), SupervisorStateEnum.STATE_RUNNING, "state4")
	supervisor.StopAndWait(2)
	AssertEqual(t, supervisor.GetState(2), SupervisorStateEnum.STATE_EXITED, "state5")
	AssertEqual(t, argvInt2, 101, "input2 argv != 101")

	supervisor.StopAndWait(1)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state6")
	AssertEqual(t, argvInt, 100, "input argv != 100")
	AssertRoutineOne(t)
}

//测试一个实例重复启动
func TestSupervisorRepeatStart(t *testing.T) {
	supervisor := NewThreadSupervisor()

	var argvInt int
	handler := func(stopChan chan bool, input int) {
		<-stopChan
		argvInt = input
	}
	supervisor.Start(1, handler, 100)
	err := supervisor.Start(1, handler, 101)
	AssertEqual(t, err != nil, true, "start error1")
	err = supervisor.Start(1, handler, 102)
	AssertEqual(t, err != nil, true, "start error2")
	supervisor.StopAndWait(1)
	AssertEqual(t, argvInt, 100, "input argv != 100")

	supervisor.Start(1, handler, 103)
	err = supervisor.Start(1, handler, 104)
	AssertEqual(t, err != nil, true, "start error3")
	err = supervisor.Start(1, handler, 105)
	AssertEqual(t, err != nil, true, "start error4")
	supervisor.StopAndWait(1)
	AssertEqual(t, argvInt, 103, "input argv != 103")
	AssertRoutineOne(t)
}

//测试手动停止而且wait
func TestTestSupervisorStopAndWait(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		for {
			isBreak := false
			select {
			case <-stopChan:
				isBreak = true
			default:
				time.Sleep(time.Millisecond * 100)
				isBreak = false
			}
			if isBreak {
				break
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	supervisor.Start(1, handler, 100)
	time.Sleep(time.Millisecond * 50)
	clock1 := time.Now().UnixNano()
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state")
	supervisor.StopAndWait(1)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state2")
	clock2 := time.Now().UnixNano()
	AssertEqual(t, (clock2-clock1)/int64(time.Millisecond) > 100, true, "wait time")
	AssertRoutineOne(t)
}

//测试手动停止不wait
func TestTestSupervisorStopAndNotWait(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		for {
			isBreak := false
			select {
			case <-stopChan:
				isBreak = true
			default:
				time.Sleep(time.Millisecond * 100)
				isBreak = false
			}
			if isBreak {
				break
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	supervisor.Start(1, handler, 100)
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state")
	supervisor.Stop(1)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_STOPPING, "state2")
	time.Sleep(time.Millisecond * 200)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state3")
	AssertRoutineOne(t)
}

//测试abort停止
func TestTestSupervisorAbort(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		time.Sleep(time.Millisecond * 100)
	}
	var abortTaskId int
	supervisor.SetAbortHandler(func(taskId int) {
		AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_EXITED, "state3")
		abortTaskId = taskId
	})
	supervisor.Start(10001, handler, 100)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_RUNNING, "state")
	time.Sleep(time.Millisecond * 50)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_RUNNING, "state2")
	time.Sleep(time.Millisecond * 100)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_EXITED, "state3")
	AssertEqual(t, abortTaskId, 10001, "abort TaskId")
	AssertRoutineOne(t)
}

//测试abort停止后手动停止
func TestTestSupervisorAbortAndStop(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		time.Sleep(time.Millisecond * 100)
	}
	supervisor.Start(1, handler, 100)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state")
	time.Sleep(time.Millisecond * 150)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state2")
	err := supervisor.StopAndWait(1)
	AssertEqual(t, err != nil, true, "stop should fail")
	AssertRoutineOne(t)
}

//测试停止后abort
func TestTestSupervisorStopAndAbort(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		time.Sleep(time.Millisecond * 100)
	}
	var abortTaskId int
	supervisor.SetAbortHandler(func(taskId int) {
		abortTaskId = taskId
	})
	supervisor.Start(1, handler, 100)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state")
	supervisor.StopAndWait(1)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state2")
	AssertEqual(t, abortTaskId, 0, "abortTaskId should be zero")
	AssertRoutineOne(t)
}

//测试重复停止
func TestTestSupervisorRepeatStop(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		<-stopChan
	}
	supervisor.Start(1, handler, 100)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_RUNNING, "state")
	supervisor.StopAndWait(1)
	AssertEqual(t, supervisor.GetState(1), SupervisorStateEnum.STATE_EXITED, "state2")
	err := supervisor.StopAndWait(1)
	AssertEqual(t, err != nil, true, "StopAndWait should error")
	AssertRoutineOne(t)
}

//测试abort停止后重启
func TestTestSupervisorAbortAndRestart(t *testing.T) {
	supervisor := NewThreadSupervisor()
	handler := func(stopChan chan bool, input int) {
		time.Sleep(time.Millisecond * 100)
	}
	supervisor.SetAbortHandler(func(taskId int) {
		AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_EXITED, "state5")
		time.Sleep(time.Millisecond * 100)
		supervisor.Start(taskId, handler, 100)
	})
	supervisor.Start(10001, handler, 100)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_RUNNING, "state")
	time.Sleep(time.Millisecond * 150)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_EXITED, "state2")
	time.Sleep(time.Millisecond * 100)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_RUNNING, "state3")
	supervisor.StopAndWait(10001)
	AssertEqual(t, supervisor.GetState(10001), SupervisorStateEnum.STATE_EXITED, "state4")
	AssertRoutineOne(t)
}
