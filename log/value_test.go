package log

import (
	"context"
	"testing"
)

func TestValue(t *testing.T) {
	var v1 interface{}
	got := Value(context.Background(), v1)
	if got != v1 {
		t.Errorf("Value() = %v, want %v", got, v1)
	}
	var v2 Valuer = func(ctx context.Context) interface{} {
		return 3
	}
	got = Value(context.Background(), v2)
	res := got.(int)
	if res != 3 {
		t.Errorf("Value() = %v, want %v", res, 3)
	}

	kvs := make([]interface{}, 2)
	kvs = append(kvs, v1)
	kvs = append(kvs, v1)
	if containsValuer(kvs) {
		t.Errorf("containsValuer(kvs) = true, want false")
	}

	kvs = append(kvs, v2)
	kvs = append(kvs, v2)
	if !containsValuer(kvs) {
		t.Errorf("containsValuer(kvs) = false, want true")
	}
}
