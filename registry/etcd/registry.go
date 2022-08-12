package etcd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/2516319251/boosters/registry"
)

type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func New(client *clientv3.Client, opts ...Option) *Registry {
	o := &options{
		ctx:       context.Background(),
		namespace: "discovery://",
		ttl:       15 * time.Second,
		maxRetry:  5,
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Registry{
		opts:   o,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

func (r *Registry) Register(ctx context.Context, instance *registry.ServiceInstance) error {
	// 准备服务实例的键值序列化
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, instance.Name, instance.ID)
	value, err := marshal(instance)
	if err != nil {
		return err
	}

	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)

	leaseID, err := r.put(ctx, key, value)
	if err != nil {
		return err
	}

	// 心跳检测
	go r.heartbeat(r.opts.ctx, leaseID, key, value)

	return nil
}

func (r *Registry) put(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}

	log.Printf("[ETCD] put %s %s", key, value)
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}

	return grant.ID, nil
}

func (r *Registry) heartbeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.Seed(time.Now().Unix())

	for {
		if curLeaseID == 0 {
			// try to registerWithKV
			retreat := []int{}
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// prevent infinite blocking
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.put(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}

func (r *Registry) Deregister(ctx context.Context, instance *registry.ServiceInstance) error {
	// 停止租约
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()

	// 删除服务实例
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, instance.Name, instance.ID)
	log.Printf("[ETCD] delete %s", key)
	_, err := r.client.Delete(ctx, key)

	return err
}
