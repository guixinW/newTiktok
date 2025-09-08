package distributed_lock

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// DistributedLock 封装了 etcd 的锁
type DistributedLock struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

// NewDistributedLock 创建一个新的分布式锁实例
// lockKey: 锁的键，通常基于资源来命名，例如 "lock/product/123"
// ttl: 锁的租约时间（秒），防止节点崩溃导致死锁
func NewDistributedLock(etcdClient *clientv3.Client, lockKey string, ttl int) (*DistributedLock, error) {
	session, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd session: %w", err)
	}
	mutex := concurrency.NewMutex(session, lockKey)
	return &DistributedLock{
		session: session,
		mutex:   mutex,
	}, nil
}

// TryLock 尝试获取锁（非阻塞或带超时）
func (l *DistributedLock) TryLock(ctx context.Context) error {
	return l.mutex.TryLock(ctx)
}

// Lock 尝试获取锁（阻塞或带超时）
func (l *DistributedLock) Lock(ctx context.Context) error {
	return l.mutex.Lock(ctx)
}

// Unlock 释放锁
func (l *DistributedLock) Unlock(ctx context.Context) error {
	defer l.session.Close()
	return l.mutex.Unlock(ctx)
}
