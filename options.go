package graceful

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Option func(g *Graceful)

type Listen func() os.Signal

func DefaultListen() os.Signal {
	var sigCh = make(chan os.Signal)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	return <-sigCh
}

// WithListenFunc 自定义信号监听方法
func WithListenFunc(f Listen) Option {
	return func(g *Graceful) {
		g.listen = f
	}
}

// WithListenLog 自定义日志输出
func WithListenLog(infoLog, errLog *log.Logger, callDepth ...int) Option {
	if len(callDepth) == 0 {
		callDepth = append(callDepth, 3)
	}
	return func(g *Graceful) {
		g.logger = &logger{
			callDepth: callDepth[0],
			infoLog:   infoLog,
			errorLog:  errLog,
		}
	}
}
