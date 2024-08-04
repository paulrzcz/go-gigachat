package gigachat

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Token(t *testing.T) {
	token := new(Token)
	token.Set("test", time.Now().Add(time.Millisecond*100))

	assert.True(t, token.Active())
	assert.Equal(t, "test", token.Get())

	time.Sleep(time.Millisecond * 150)

	assert.False(t, token.Active())
	assert.Equal(t, "", token.Get())
}
