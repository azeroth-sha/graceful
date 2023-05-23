# graceful

优雅启停服务

### 说明
- 用于优雅的启停多个服务，对多服务应用非常有效
- 以下只需要简单的封装即可实现多服务兼容
- http.Server
- [gin](https://github.com/gin-gonic/gin)
- [rpcx](https://github.com/smallnest/rpcx)
- [gnet](https://github.com/panjf2000/gnet)
- [cron](https://github.com/robfig/cron)
- 更多

### 用法、示例

- main.go
```go
package main

import (
	"github.com/azeroth-sha/graceful"
	"log"
)

func init() {
	_ = graceful.Add(`gin`, new(ginSvr))
	_ = graceful.Add("cron", new(cronSvr))
}

func main() {
	if err := graceful.Run(); err != nil {
		log.Print(err)
	}
	if err := graceful.Stop(); err != nil {
		log.Print(err)
	}
}
```

- gin.go
```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
)

type ginSvr struct {
	once    sync.Once
	running int32
	r       *gin.Engine
	svr     *http.Server
}

func (g *ginSvr) init() {
	g.r = gin.New()
	g.svr = &http.Server{
		Addr:    ":8080",
		Handler: g.r,
	}
}

func (g *ginSvr) Service() error {
	if atomic.SwapInt32(&g.running, 1) != 0 {
		return nil
	}
	g.once.Do(g.init)
	return g.svr.ListenAndServe()
}

func (g *ginSvr) Shutdown(ctx context.Context) error {
	if atomic.SwapInt32(&g.running, 0) != 1 {
		return nil
	}
	return g.svr.Shutdown(ctx)
}
```

- cron.go
```go
package main

import (
	"context"
	cron "github.com/robfig/cron/v3"
	"sync"
	"sync/atomic"
)

type cronSvr struct {
	running int32
	once    sync.Once
	svr     *cron.Cron
}

func (e *cronSvr) init() {
	e.svr = cron.New()
}

func (e *cronSvr) Service() error {
	if atomic.SwapInt32(&e.running, 1) != 0 {
		return nil
	}
	e.once.Do(e.init)
	e.svr.Run()
	return nil
}

func (e *cronSvr) Shutdown(ctx context.Context) error {
	if atomic.SwapInt32(&e.running, 0) != 1 {
		return nil
	}
	e.svr.Stop()
	return nil
}
```