package main

import (
	"errors"
	"fmt"
)

const maxTrxDepth = 100

var (
	ErrNoValidTrx       = errors.New("there is no valid transaction existed")
	ErrMaxDepthExceeded = errors.New("max transaction depth exceeded")
)

type Session struct {
	topTrx  *transaction
	trxSize int
}

func NewSession() *Session {
	return &Session{
		topTrx:  newTransaction(),
		trxSize: 1,
	}
}

func (_s *Session) Begin() error {
	if _s.trxSize == maxTrxDepth {
		return ErrMaxDepthExceeded
	}
	trx := newTransaction()
	trx.Push(_s.topTrx)
	_s.topTrx, _s.trxSize = trx, _s.trxSize+1
	return nil
}

func (_s *Session) Rollback() error {
	if _s.trxSize == 1 {
		return ErrNoValidTrx
	}
	_s.topTrx = _s.topTrx.Next()
	_s.trxSize -= 1
	return nil
}

func (_s *Session) Commit() error {
	if _s.trxSize == 1 {
		return ErrNoValidTrx
	}
	next := _s.topTrx.Next()
	for k, v := range _s.topTrx.store {
		next.Set(k, v)
	}
	_s.topTrx, _s.trxSize = next, _s.trxSize-1
	return nil
}

// Count function is operated in-memory, the iteration is so fast,
// based on the requirements, no need to add a reversed map to archive O(1) complexity.
func (_s *Session) Count(value string) int {
	result, current := 0, _s.topTrx
	for current != nil {
		for _, v := range current.store {
			if value == v {
				result += 1
			}
		}
		current = current.Next()
	}
	return result
}

func (_s *Session) Delete(key string) {
	_s.topTrx.Delete(key)
}

func (_s *Session) Set(key, value string) {
	_s.topTrx.Set(key, value)
}

func (_s *Session) Get(key string) (string, bool) {
	next := _s.topTrx
	for next != nil {
		v, exist := _s.topTrx.Get(key)
		if exist {
			return v, true
		}
		next = next.Next()
	}
	return "", false
}

func (_s *Session) MustGet(key string) string {
	v, exist := _s.Get(key)
	if !exist {
		panic(fmt.Sprintf("the value of %s should be provided", key))
	}
	return v
}
