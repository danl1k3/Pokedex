package pokecache

import (
	"bytes"
	"testing"
	"time"
	"sync"
)

func TestCache_AddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		value []byte
	}{
		{
			key: "https://example.com",
			value: []byte("testdata"),
		},
		{
			key: "Https://example.com/path",
			value: []byte("moretestdata"),
		},
	}

	cache := NewCache(interval)

	for _, c := range cases {
		cache.Add(c.key, c.value)
		actual, ok := cache.Get(c.key)
		if !ok {
			t.Fatalf("Error while getting key: %s", c.key)
		}
		if !bytes.Equal(actual, c.value) {
			t.Errorf("Got: %v, want: %v", actual, c.value)
		}
	}
}

func TestCache_GetMissing(t *testing.T) {
	cache := NewCache(5 * time.Second)

	actual, ok := cache.Get("nonexistent")
	if ok {
		t.Errorf("Retrived non-existent data!")
	}
	if actual != nil {
		t.Errorf("Retrived non-existent value! Got: %v", actual)
	}
}

func TestCache_Overwrite(t *testing.T) {
	const interval = 40 * time.Millisecond
	cache := NewCache(interval)
	key := "https://example.com"

	cache.Add(key, []byte("old"))

	time.Sleep(25 * time.Millisecond)
	cache.Add(key, []byte("new"))

	time.Sleep(25 * time.Millisecond)

	actual, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get %v fails to retrive from cache", key)
	}
	if !bytes.Equal(actual, []byte("new")){
		t.Errorf("Got: %s, want %s", actual, "new")
	}
}
func TestCache_ReapLoop(t *testing.T) {
	const baseTime = 10 * time.Millisecond
	const waitTime = baseTime + 15*time.Millisecond

	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	if _, ok := cache.Get("https://example.com"); !ok {
		t.Fatalf("expected to find key before expiration")
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get("https://example.com"); ok {
		t.Fatalf("expected key to be evicted after interval")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			cache.Add("shared-key", []byte("value"))
		}()
		go func() {
			defer wg.Done()
			cache.Get("shared-key")
		}()
	}

	wg.Wait()
}
