package graceful

import (
	"context"
	"runtime/debug"
	"sync/atomic"
)

type service struct {
	g       *Graceful
	running int32
	name    string
	svr     Service
}

func (s *service) run() {
	defer func() {
		if err := recover(); err != nil {
			s.g.log().Errorf("%v\r\n%s", debug.Stack())
		}
	}()
	if err := s.svr.Service(); err != nil {
		s.g.log().Error(err)
	}
}

func (s *service) start() {
	if atomic.SwapInt32(&s.running, 1) != 0 {
		return
	}
	go func() {
		s.g.log().Infof("%s: running", s.name)
		for atomic.LoadInt32(&s.running) == 1 {
			s.run()
		}
	}()
}

func (s *service) stop() {
	if atomic.SwapInt32(&s.running, 0) != 1 {
		return
	}
	defer s.g.log().Infof("%s: stopped", s.name)
	if err := s.svr.Shutdown(context.Background()); err != nil {
		s.g.log().Error(err)
	}
}
