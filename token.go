package gigachat

import (
	"sync"
	"time"
)

type Token struct {
	expiresAt time.Time
	value     string
	mx        sync.RWMutex
}

func (t *Token) Active() bool {
	t.mx.RLock()
	defer t.mx.RUnlock()

	return time.Now().Before(t.expiresAt)
}

func (t *Token) Get() string {
	t.mx.RLock()
	defer t.mx.RUnlock()

	if t.Active() {
		return t.value
	}

	return ""
}

func (t *Token) Set(value string, expiresAt time.Time) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.value = value
	t.expiresAt = expiresAt
}
