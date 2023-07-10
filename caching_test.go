package caching_test

import (
	"testing"
	"time"

	"github.com/g3offrey/caching"
)

func getTimestamp() int64 {
	return time.Now().UnixMilli()
}

func TestCacheRemember(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	v := c.Remember("timestamp", getTimestamp)
	time.Sleep(50 * time.Millisecond)
	v2 := c.Remember("timestamp", getTimestamp)

	if v != v2 {
		t.Errorf("Expected %d to be equal to %d", v, v2)
	}
}

func TestCacheRememberRecompute(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	v := c.Remember("timestamp", getTimestamp)
	time.Sleep(200 * time.Millisecond)
	v2 := c.Remember("timestamp", getTimestamp)

	if v == v2 {
		t.Errorf("Expected %d to be different from %d", v, v2)
	}
}

func TestCacheGetStaleThenRecompute(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	v := c.GetStaleThenRecompute("timestamp", getTimestamp)
	time.Sleep(200 * time.Millisecond)
	v2 := c.GetStaleThenRecompute("timestamp", getTimestamp)
	time.Sleep(50 * time.Millisecond)
	v3 := c.GetStaleThenRecompute("timestamp", getTimestamp)

	if v != v2 {
		t.Errorf("Expected %d to be equal to %d", v, v2)
	}
	if v2 == v3 {
		t.Errorf("Expected %d to be different from %d", v2, v3)
	}
}

func TestCacheGet(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	_ = c.Remember("timestamp", func() int64 { return 10 })
	v, expired, err := c.Get("timestamp")

	if v != 10 {
		t.Errorf("Expected %d to be equal to %d", v, 10)
	}
	if expired {
		t.Errorf("Expected %t to be equal to %t", expired, false)
	}
	if err != nil {
		t.Errorf("Expected %s to be nil", err)
	}
}

func TestCacheGetWithExpiredValue(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	_ = c.Remember("timestamp", func() int64 { return 10 })
	time.Sleep(200 * time.Millisecond)
	v, expired, err := c.Get("timestamp")

	if v != 10 {
		t.Errorf("Expected %d to be equal to %d", v, 10)
	}
	if !expired {
		t.Errorf("Expected %t to be equal to %t", expired, true)
	}
	if err != nil {
		t.Errorf("Expected %s to be nil", err)
	}
}

func TestCacheGetWithNonExistingValue(t *testing.T) {
	c := caching.New[int64](100 * time.Millisecond)

	v, expired, err := c.Get("timestamp")

	if v != 0 {
		t.Errorf("Expected %d to be equal to %d", v, 0)
	}
	if expired {
		t.Errorf("Expected %t to be equal to %t", expired, false)
	}
	if err == nil {
		t.Errorf("Expected %s to not be nil", err)
	}
}
