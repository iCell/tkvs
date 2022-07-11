package tkvs

import (
	"errors"
	"fmt"
)

const maxTrxDepth = 100

var (
	ErrMaxDepthExceeded = errors.New("max transaction depth exceeded")
	ErrNoValidTrx       = errors.New("no transaction")
)

type KvStore struct {
	topTrx   *transaction
	trxCount int
}

func NewKvStore() *KvStore {
	return &KvStore{
		topTrx:   newTransaction(),
		trxCount: 1,
	}
}

func (_kv *KvStore) Begin() error {
	if _kv.trxCount == maxTrxDepth {
		return ErrMaxDepthExceeded
	}
	trx := newTransaction()
	trx.Next = _kv.topTrx
	_kv.topTrx, _kv.trxCount = trx, _kv.trxCount+1
	return nil
}

func (_kv *KvStore) Rollback() error {
	if _kv.trxCount == 1 {
		return ErrNoValidTrx
	}
	_kv.topTrx, _kv.trxCount = _kv.topTrx.Next, _kv.trxCount-1
	return nil
}

func (_kv *KvStore) Commit() error {
	if _kv.trxCount == 1 {
		return ErrNoValidTrx
	}
	next := _kv.topTrx.Next
	for k, v := range _kv.topTrx.Kvs {
		next.Set(k, v)
	}
	_kv.topTrx, _kv.trxCount = next, _kv.trxCount-1
	return nil
}

// Count function is operated in-memory, the iteration is so fast,
// based on the requirements, no need to add a reversed map to archive O(1) complexity.
func (_kv *KvStore) Count(value string) int {
	result, visited, current := 0, make(map[string]bool), _kv.topTrx
	for current != nil {
		for key, v := range current.Kvs {
			if !visited[key] && value == v {
				result += 1
			}
			visited[key] = true
		}
		current = current.Next
	}
	return result
}

func (_kv *KvStore) Delete(key string) {
	_kv.topTrx.Delete(key)
}

func (_kv *KvStore) Set(key, value string) {
	_kv.topTrx.Set(key, value)
}

func (_kv *KvStore) Get(key string) (string, bool) {
	next := _kv.topTrx
	for next != nil {
		v, exist := next.Get(key)
		if exist {
			return v, true
		}
		next = next.Next
	}
	return "", false
}

func (_kv *KvStore) MustGet(key string) string {
	v, exist := _kv.Get(key)
	if !exist {
		panic(fmt.Sprintf("the value of %s should be provided", key))
	}
	return v
}
