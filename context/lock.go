package context

import (
	"github.com/cloudimpl/byte-os/runtime"
	"time"
)

type Lock struct {
	client    runtime.ServiceClient
	sessionId string
	key       string
}

func (l *Lock) Acquire(expireIn time.Duration) error {
	ttl := time.Now().Unix() + int64(expireIn.Seconds())

	req := runtime.AcquireLockRequest{
		Key: l.key,
		TTL: ttl,
	}

	return l.client.AcquireLock(l.sessionId, req)
}

func (l *Lock) Release() error {
	req := runtime.ReleaseLockRequest{
		Key: l.key,
	}

	return l.client.ReleaseLock(l.sessionId, req)
}
