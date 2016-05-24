package web

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"github.com/facebookgo/httpdown"
	. "github.com/fishedee/language"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var (
	didInherit = os.Getenv("LISTEN_FDS") != ""
	ppid       = os.Getppid()
)

type Grace interface {
	ListenAndServe(port int, handler http.Handler) error
}

type GraceConfig struct {
	Driver  string
	Stop    []string
	Restart []string
}

type graceImplement struct {
	runGrace      bool
	graceServer   httpdown.Server
	graceNet      *gracenet.Net
	stopSignal    map[os.Signal]bool
	restartSignal map[os.Signal]bool
	waitGroup     sync.WaitGroup
}

func NewGrace(config GraceConfig) (Grace, error) {
	var goGrace bool
	if config.Driver != "" {
		goGrace = true
	} else {
		goGrace = false
	}
	stopSignal := map[os.Signal]bool{}
	for _, singleStop := range config.Stop {
		signal, err := stringToSignal(singleStop)
		if err != nil {
			return nil, err
		}
		stopSignal[signal] = true
	}
	restartSignal := map[os.Signal]bool{}
	for _, singleRestart := range config.Restart {
		signal, err := stringToSignal(singleRestart)
		if err != nil {
			return nil, err
		}
		restartSignal[signal] = true
	}
	return &graceImplement{
		runGrace:      goGrace,
		stopSignal:    stopSignal,
		restartSignal: restartSignal,
	}, nil
}

func NewGraceFromConfig(configName string) (Grace, error) {
	gracedirver := globalBasic.Config.GetString(configName + "driver")
	gracestopStr := globalBasic.Config.GetString(configName + "stop")
	gracerestartStr := globalBasic.Config.GetString(configName + "restart")

	gracestop := Explode(gracestopStr, ",")
	gracerestart := Explode(gracerestartStr, ",")
	config := GraceConfig{
		Driver:  gracedirver,
		Stop:    gracestop,
		Restart: gracerestart,
	}
	return NewGrace(config)
}

func stringToSignal(signalStr string) (os.Signal, error) {
	result := map[string]os.Signal{
		"TERM": syscall.SIGTERM,
		"INT":  syscall.SIGINT,
		"HUP":  syscall.SIGHUP,
		"USR1": syscall.SIGUSR1,
		"USR2": syscall.SIGUSR2,
	}
	target, isExist := result[signalStr]
	if isExist == false {
		return nil, fmt.Errorf("invalid signal %v", signalStr)
	} else {
		return target, nil
	}
}

func (this *graceImplement) signalHandler(errorEvent chan error) {
	ch := make(chan os.Signal, 10)

	allSignal := []os.Signal{}
	for signal, _ := range this.stopSignal {
		allSignal = append(allSignal, signal)
	}
	for signal, _ := range this.restartSignal {
		allSignal = append(allSignal, signal)
	}
	signal.Notify(ch, allSignal...)

	for {
		sig := <-ch
		if _, isStopSignal := this.stopSignal[sig]; isStopSignal {
			globalBasic.Log.Debug("Receive Stop Signal")
			//优雅关闭
			go func() {
				defer this.waitGroup.Done()
				if err := this.graceServer.Stop(); err != nil {
					errorEvent <- err
				}
			}()
		} else {
			globalBasic.Log.Debug("Receive Restart Signal")
			//优雅重启
			if _, err := this.graceNet.StartProcess(); err != nil {
				errorEvent <- err
			}
		}
	}
}

func (this *graceImplement) waitHandler(errorEvent chan error) {
	defer this.waitGroup.Done()
	if err := this.graceServer.Wait(); err != nil {
		errorEvent <- err
	}
}

func (this *graceImplement) waitServerStop(errorEvent chan error) {
	this.waitGroup.Add(2)
	go this.signalHandler(errorEvent)
	go this.waitHandler(errorEvent)
	this.waitGroup.Wait()
}

func (this *graceImplement) listenAndServeGrace(httpPort string, handler http.Handler) error {
	//倾听端口
	this.graceNet = &gracenet.Net{}
	listener, err := this.graceNet.Listen("tcp", httpPort)
	if err != nil {
		return err
	}

	if didInherit {
		if ppid == 1 {
			globalBasic.Log.Debug("Listening on init activated %v", httpPort)
		} else {
			globalBasic.Log.Debug("Graceful handoff of %v with new pid %v and old pid %v", httpPort, os.Getpid(), ppid)
		}
	} else {
		globalBasic.Log.Debug("Serving %s with pid %d", httpPort, os.Getpid())
	}

	//对外服务
	httpServer := &httpdown.HTTP{}
	this.graceServer = httpServer.Serve(&http.Server{
		Addr:    httpPort,
		Handler: handler,
	}, listener)

	//关闭父级进程
	if didInherit && ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to close parent: %s", err)
		}
	}

	//等待服务器结束
	errorEvent := make(chan error)
	waitEvent := make(chan bool)
	go func() {
		defer close(waitEvent)
		this.waitServerStop(errorEvent)
	}()

	select {
	case err := <-errorEvent:
		if err == nil {
			panic("unexpected nil error")
		}
		return err
	case <-waitEvent:
		globalBasic.Log.Debug("Exiting pid %v.", os.Getpid())
		return nil
	}
}

func (this *graceImplement) ListenAndServe(httpPort int, handler http.Handler) error {
	if this.runGrace == false {
		return http.ListenAndServe(":"+strconv.Itoa(httpPort), handler)
	} else {
		return this.listenAndServeGrace(":"+strconv.Itoa(httpPort), handler)
	}
}
