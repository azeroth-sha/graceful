package graceful

import (
	"context"
	"errors"
	"sync/atomic"
)

var (
	ErrRunning = errors.New(`graceful is running`)
	ErrStopped = errors.New(`graceful is stopped`)
	g          = New()
)

// Service 服务方法
type Service interface {
	Service() error
	Shutdown(ctx context.Context) error
}

type Graceful struct {
	running  int32
	listen   Listen
	logger   *logger
	services map[string]*service
}

func (g *Graceful) log() *logger {
	if g.logger != nil {
		return g.logger
	}
	return defaultLog
}

// Add 添加服务
func (g *Graceful) Add(name string, svr Service) error {
	if atomic.LoadInt32(&g.running) == 1 {
		return ErrRunning
	}
	g.services[name] = &service{
		g:       g,
		running: 0,
		name:    name,
		svr:     svr,
	}
	return nil
}

// Run 运行所有服务
func (g *Graceful) Run() error {
	if atomic.SwapInt32(&g.running, 1) != 0 {
		return ErrRunning
	}
	for _, svr := range g.services {
		svr.start()
	}
	g.log().Infof("received signal: %s", g.listen())
	return nil
}

// Stop 停止所有服务
func (g *Graceful) Stop() error {
	if atomic.SwapInt32(&g.running, 0) != 1 {
		return ErrStopped
	}
	for _, svr := range g.services {
		svr.stop()
	}
	return nil
}

// New 创建服务维护对象
func New(opts ...Option) *Graceful {
	graceful := &Graceful{
		running:  0,
		listen:   DefaultListen,
		logger:   defaultLog,
		services: make(map[string]*service),
	}
	for _, opt := range opts {
		opt(graceful)
	}
	return graceful
}

// Add 添加服务
func Add(name string, svr Service) error {
	return g.Add(name, svr)
}

// Run 运行所有服务
func Run() error {
	return g.Run()
}

// Stop 停止所有服务
func Stop() error {
	return g.Stop()
}
