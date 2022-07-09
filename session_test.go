package main

import "testing"

func TestSession(t *testing.T) {
	t.Run("test basic set/get/delete without transaction", func(t *testing.T) {
		session := NewSession()
		session.Set("key1", "value1")
		session.Set("key2", "value2")
		v1, exist := session.Get("key1")
		assert(t, true, exist)
		assert(t, "value1", v1)
		session.Delete("key2")
		v2, exist := session.Get("key2")
		assert(t, false, exist)
		assert(t, "", v2)
		assert(t, 1, session.Count("value1"))
		assert(t, 0, session.Count("value2"))
	})

	t.Run("test set/get then start a transaction and commit", func(t *testing.T) {
		session := NewSession()
		session.Set("key1", "value1")
		session.Set("key2", "value2")
		assert(t, "value1", session.MustGet("key1"))
		assert(t, "value2", session.MustGet("key2"))

		session.Begin()
		session.Set("key1", "new_value1")
		session.Set("key3", "new_value3")
		err := session.Commit()
		shouldBeNoError(t, err)
		assert(t, "new_value1", session.MustGet("key1"))
		assert(t, "new_value3", session.MustGet("key3"))
		assert(t, "value2", session.MustGet("key2"))
	})

	t.Run("test set/get then start a transaction and rollback", func(t *testing.T) {
		session := NewSession()
		session.Set("key1", "value1")
		session.Set("key2", "value2")
		assert(t, "value1", session.MustGet("key1"))
		assert(t, "value2", session.MustGet("key2"))

		session.Begin()
		session.Set("key1", "new_value1")
		err := session.Rollback()
		shouldBeNoError(t, err)
		assert(t, "value1", session.MustGet("key1"))
		assert(t, "value2", session.MustGet("key2"))
	})

	t.Run("test start a transaction with nested transaction and commit", func(t *testing.T) {
		session := NewSession()

		session.Begin()
		session.Set("key1", "value1")

		session.Begin()
		session.Set("inner_key1", "inner_value1")
		session.Set("inner_key2", "inner_value2")

		session.Begin()
		session.Set("inner_key3", "inner_value3")
		session.Commit()
		session.Commit()

		session.Commit()

		assert(t, "value1", session.MustGet("key1"))
		assert(t, "inner_value1", session.MustGet("inner_key1"))
		assert(t, "inner_value2", session.MustGet("inner_key2"))
		assert(t, "inner_value3", session.MustGet("inner_key3"))
	})

	t.Run("test start a transaction with nested transaction and rollback one of the nested transaction", func(t *testing.T) {
		session := NewSession()
		session.Set("key", "value")

		session.Begin()
		session.Set("key1", "value1")

		session.Begin()
		session.Set("key", "new_value")
		session.Set("inner_key1", "inner_value1")
		session.Set("inner_key2", "inner_value2")

		session.Begin()
		session.Set("inner_key3", "inner_value3")
		session.Rollback()

		session.Commit()

		session.Commit()

		assert(t, "new_value", session.MustGet("key"))
		assert(t, "value1", session.MustGet("key1"))
		assert(t, "inner_value1", session.MustGet("inner_key1"))
		assert(t, "inner_value2", session.MustGet("inner_key2"))
		_, exist := session.Get("inner_key3")
		assert(t, false, exist)
	})

	t.Run("test start a transaction with nested transaction and commit, to check the count operation's result", func(t *testing.T) {
		session := NewSession()
		session.Set("key", "value")

		session.Begin()
		session.Set("key1", "value1")

		session.Begin()
		session.Set("key", "new_value")
		session.Set("inner_key1", "value1")
		session.Set("inner_key2", "inner_value2")
		session.Commit()

		session.Commit()

		assert(t, 1, session.Count("new_value"))
		assert(t, 2, session.Count("value1"))
		assert(t, 1, session.Count("inner_value2"))
		assert(t, 0, session.Count("value"))
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
