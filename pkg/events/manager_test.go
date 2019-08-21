package events

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPubSub(t *testing.T) {
	ctx := context.TODO()
	m := NewManager(ctx, 32)
	called := false
	_, err := m.Subscribe("test", func(interface{}) error {
		called = true
		return nil
	})
	require.Nil(t, err)
	err = m.Publish("test", nil)
	require.Nil(t, err)
	require.True(t, called)
}

func TestConcurrency(t *testing.T) {
	ctx := context.TODO()
	m := NewManager(ctx, 32)
	var counter uint64
	for i := 0; i < 1024; i++ {
		_, err := m.Subscribe("test", func(interface{}) error {
			atomic.AddUint64(&counter, 1)
			return nil
		})
		require.Nil(t, err)
	}
	err := m.Publish("test", nil)
	require.Nil(t, err)
	require.Equal(t, uint64(1024), counter)
}

func TestUnsubscribe(t *testing.T) {
	ctx := context.TODO()
	m := NewManager(ctx, 32)
	called := false
	id, err := m.Subscribe("test", func(interface{}) error {
		called = true
		return nil
	})
	require.Nil(t, err)
	err = m.Unsubscribe("test", id)
	require.Nil(t, err)
	err = m.Publish("test", nil)
	require.Nil(t, err)
	require.False(t, called)
}

func TestErrorHandling(t *testing.T) {
	ctx := context.TODO()
	m := NewManager(ctx, 32)
	preparedError := errors.New("this is the error")
	_, err := m.Subscribe("test", func(interface{}) error {
		return preparedError
	})
	require.Nil(t, err)
	err = m.Publish("test", nil)
	require.NotNil(t, err)
	require.Equal(t, err.(MultiError)[0], errors.New("this is the error"))
}

func BenchmarkPublishing(b *testing.B) {
	ctx := context.TODO()
	m := NewManager(ctx, 32)
	m.Subscribe("test", func(interface{}) error { return nil })
	for n := 0; n < b.N; n++ {
		m.Publish("test", nil)
	}
}
