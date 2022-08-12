package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/2516319251/boosters/internal/endpoint"
	"github.com/2516319251/boosters/registry"
)

type discoveryResolver struct {
	watcher registry.Watcher
	conn    resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure         bool
	debugLogDisabled bool
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		instance, err := r.watcher.Next()
		if err != nil {
			// 如果是取消上下文
			if errors.Is(err, context.Canceled) {
				return
			}
			// 打印错误，并继续
			log.Printf("[RESOLVER] failed to watch discovery endpoint: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(instance)
	}
}

func (r *discoveryResolver) update(instances []*registry.ServiceInstance) {
	addrs := make([]resolver.Address, 0)
	endpoints := make(map[string]struct{})

	for _, instance := range instances {
		ept, err := endpoint.ParseEndpoint(instance.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		if err != nil {
			log.Printf("[RESOLVER] Failed to parse discovery endpoint: %v", err)
			continue
		}

		if ept == "" {
			continue
		}

		// filter redundant endpoints
		if _, ok := endpoints[ept]; ok {
			continue
		}

		endpoints[ept] = struct{}{}
		addr := resolver.Address{
			ServerName: instance.Name,
			Attributes: parseAttributes(instance.Metadata),
			Addr:       ept,
		}
		addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", instance)
		addrs = append(addrs, addr)
	}

	if len(addrs) == 0 {
		log.Printf("[RESOLVER] Zero endpoint found,refused to write, instances: %v", instances)
		return
	}

	err := r.conn.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Printf("[RESOLVER] failed to update state: %s", err)
	}

	if !r.debugLogDisabled {
		b, _ := json.Marshal(instances)
		log.Printf("[RESOLVER] update instances: %s", b)
	}
}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}

func (r *discoveryResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.watcher.Stop()
	if err != nil {
		log.Printf("[RESOLVER] failed to watch top: %s", err)
	}
}
