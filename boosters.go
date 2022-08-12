package boosters

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/2516319251/boosters/registry"
)

type Boosters struct {
	opts options

	ctx    context.Context
	cancel context.CancelFunc

	mutex    sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *Boosters {
	o := options{
		ctx:  context.Background(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &Boosters{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (boosters *Boosters) Run() error {
	wg := sync.WaitGroup{}
	eg, ctx := errgroup.WithContext(boosters.ctx)

	// 服务注册
	if err := boosters.registrar(ctx); err != nil {
		return err
	}

	// 处理服务生命周期
	for _, server := range boosters.opts.servers {
		srv := server

		// 接收停止服务信号
		eg.Go(func() error {
			<-ctx.Done()
			return srv.Stop(ctx)
		})

		// 启动服务
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()

	// 接收信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, boosters.opts.sigs...)

	// 停止服务协程
	eg.Go(func() error {
		for {
			select {
			// 上下文取消
			case <-ctx.Done():
				return ctx.Err()

			// 信号
			case <-quit:
				return boosters.Stop()
			}
		}
	})

	// 阻塞在此等待服务停止
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (boosters *Boosters) registrar(ctx context.Context) error {
	// 注意是否会发生重入锁
	boosters.mutex.Lock()
	defer boosters.mutex.Unlock()

	// 如果使用服务注册
	if boosters.opts.registrar != nil {

		// 获取服务实例
		instance, err := boosters.instances()
		if err != nil {
			return err
		}

		// 停止服务时服务注销需要使用服务实例
		boosters.instance = instance

		// 服务注册
		log.Printf("[REGISTRY] register service instance")
		return boosters.opts.registrar.Register(ctx, instance)
	}

	return nil
}

func (boosters *Boosters) instances() (*registry.ServiceInstance, error) {
	log.Printf("[REGISTRY] build service instance")

	// 获取服务地址
	endpoints := make([]string, 0, len(boosters.opts.endpoints))
	for _, endpoint := range boosters.opts.endpoints {
		endpoints = append(endpoints, endpoint.String())
	}

	// 如果未设置服务端地址
	if len(endpoints) == 0 {
		return nil, errors.New("build service instance with empty endpoints")
	}

	// 服务实例信息
	return &registry.ServiceInstance{
		ID:        boosters.ID(),
		Name:      boosters.Name(),
		Version:   boosters.Version(),
		Endpoints: endpoints,
	}, nil
}

func (boosters *Boosters) Stop() error {
	// 注意是否会发生重入锁
	boosters.mutex.Lock()
	defer boosters.mutex.Unlock()

	// 如果使用了服务注册并且存在服务实例
	if boosters.opts.registrar != nil && boosters.instance != nil {
		// 服务注销
		log.Printf("[REGISTRY] deregister service instance")
		if err := boosters.opts.registrar.Deregister(boosters.ctx, boosters.instance); err != nil {
			return err
		}
	}

	// 取消上下文以停止服务
	if boosters.cancel != nil {
		boosters.cancel()
	}

	return nil
}
