package pokecache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("test1"),
		},
		{
			key: "https://example.com/path",
			val: []byte("test2"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			cache := NewCache(ctx, interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected value for key %s not found", c.key)
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value %s, got %s", c.val, val)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(ctx, baseTime)
	cache.Add("https://example.com", []byte("test"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected key not found")
		return
	}

	time.Sleep(waitTime)
	_, ok = cache.Get("https://example.com")
	if ok {
		t.Error("expected key found after cache expiration")
	}
}
