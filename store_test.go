package main

import "testing"

func TestKvStore(t *testing.T) {
	t.Run("test basic set/get/delete without transaction", func(t *testing.T) {
		kvStore := NewKvStore()
		kvStore.Set("key1", "value1")
		kvStore.Set("key2", "value2")
		v1, exist := kvStore.Get("key1")
		assert(t, true, exist)
		assert(t, "value1", v1)
		kvStore.Delete("key2")
		v2, exist := kvStore.Get("key2")
		assert(t, false, exist)
		assert(t, "", v2)
		assert(t, 1, kvStore.Count("value1"))
		assert(t, 0, kvStore.Count("value2"))
	})

	t.Run("test set/get then start a transaction and commit", func(t *testing.T) {
		kvStore := NewKvStore()
		kvStore.Set("key1", "value1")
		kvStore.Set("key2", "value2")
		assert(t, "value1", kvStore.MustGet("key1"))
		assert(t, "value2", kvStore.MustGet("key2"))

		kvStore.Begin()
		kvStore.Set("key1", "new_value1")
		kvStore.Set("key3", "new_value3")
		err := kvStore.Commit()
		shouldBeNoError(t, err)
		assert(t, "new_value1", kvStore.MustGet("key1"))
		assert(t, "new_value3", kvStore.MustGet("key3"))
		assert(t, "value2", kvStore.MustGet("key2"))
	})

	t.Run("test set/get then start a transaction and rollback", func(t *testing.T) {
		kvStore := NewKvStore()
		kvStore.Set("key1", "value1")
		kvStore.Set("key2", "value2")
		assert(t, "value1", kvStore.MustGet("key1"))
		assert(t, "value2", kvStore.MustGet("key2"))

		kvStore.Begin()
		kvStore.Set("key1", "new_value1")
		err := kvStore.Rollback()
		shouldBeNoError(t, err)
		assert(t, "value1", kvStore.MustGet("key1"))
		assert(t, "value2", kvStore.MustGet("key2"))
	})

	t.Run("test start a transaction with nested transaction and commit", func(t *testing.T) {
		kvStore := NewKvStore()

		kvStore.Begin()
		kvStore.Set("key1", "value1")

		kvStore.Begin()
		kvStore.Set("inner_key1", "inner_value1")
		kvStore.Set("inner_key2", "inner_value2")

		kvStore.Begin()
		kvStore.Set("inner_key3", "inner_value3")
		kvStore.Commit()
		kvStore.Commit()

		kvStore.Commit()

		assert(t, "value1", kvStore.MustGet("key1"))
		assert(t, "inner_value1", kvStore.MustGet("inner_key1"))
		assert(t, "inner_value2", kvStore.MustGet("inner_key2"))
		assert(t, "inner_value3", kvStore.MustGet("inner_key3"))
	})

	t.Run("test start a transaction with nested transaction and rollback one of the nested transaction", func(t *testing.T) {
		kvStore := NewKvStore()
		kvStore.Set("key", "value")

		kvStore.Begin()
		kvStore.Set("key1", "value1")

		kvStore.Begin()
		kvStore.Set("key", "new_value")
		kvStore.Set("inner_key1", "inner_value1")
		kvStore.Set("inner_key2", "inner_value2")

		kvStore.Begin()
		kvStore.Set("inner_key3", "inner_value3")
		kvStore.Rollback()

		kvStore.Commit()

		kvStore.Commit()

		assert(t, "new_value", kvStore.MustGet("key"))
		assert(t, "value1", kvStore.MustGet("key1"))
		assert(t, "inner_value1", kvStore.MustGet("inner_key1"))
		assert(t, "inner_value2", kvStore.MustGet("inner_key2"))
		_, exist := kvStore.Get("inner_key3")
		assert(t, false, exist)
	})

	t.Run("test start a transaction with nested transaction and commit, to check the count operation's result", func(t *testing.T) {
		kvStore := NewKvStore()
		kvStore.Set("key", "value")

		kvStore.Begin()
		kvStore.Set("key1", "value1")

		kvStore.Begin()
		kvStore.Set("key", "new_value")
		kvStore.Set("inner_key1", "value1")
		kvStore.Set("inner_key2", "inner_value2")
		kvStore.Commit()

		kvStore.Commit()

		assert(t, 1, kvStore.Count("new_value"))
		assert(t, 2, kvStore.Count("value1"))
		assert(t, 1, kvStore.Count("inner_value2"))
		assert(t, 0, kvStore.Count("value"))
	})
}

func shouldBeNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expect %v, but actual is %v", expected, actual)
		return
	}
}
