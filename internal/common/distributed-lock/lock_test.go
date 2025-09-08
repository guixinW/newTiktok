package distributed_lock

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	var port string
	if port = os.Getenv("ETCD_ENDPOINTS"); port == "" {
		port = "127.0.0.1:2379"
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(port, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cli.Sync(ctx)
	if err != nil {
		t.Errorf("与 etcd 集群的连接测试失败: %v", err)
	}
}

func TestLock(t *testing.T) {
	var port string
	if port = os.Getenv("ETCD_ENDPOINTS"); port == "" {
		port = "127.0.0.1:2379"
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(port, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	s1, _ := concurrency.NewSession(cli, concurrency.WithTTL(10))
	defer s1.Close()
	s2, _ := concurrency.NewSession(cli)
	defer s2.Close()
	l1 := concurrency.NewMutex(s1, "/distributed-lock")
	l2 := concurrency.NewMutex(s2, "/distributed-lock")
	ctx, cancel := context.WithTimeout(context.Background(), 5001*time.Millisecond)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := l1.Lock(ctx); err != nil {
			t.Error(err)
		}
		time.Sleep(3 * time.Second)
		t.Log("session1 lock")
		err := l1.Unlock(ctx)
		t.Log("session1 unlock")
		if err != nil {
			t.Error(err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := l2.Lock(ctx); err != nil {
			t.Error(err)
		}
		t.Log("session2 lock")
		time.Sleep(2 * time.Second)
		t.Log("session2 unlock")
		err := l2.Unlock(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()
}
