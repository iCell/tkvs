package main

type transaction struct {
	store map[string]string
	next  *transaction
}

func newTransaction() *transaction {
	return &transaction{store: make(map[string]string)}
}

func (_trx *transaction) Set(k, v string) {
	_trx.store[k] = v
}

func (_trx *transaction) Get(k string) (string, bool) {
	v, exist := _trx.store[k]
	return v, exist
}

func (_trx *transaction) Delete(k string) {
	delete(_trx.store, k)
}

func (_trx *transaction) Next() *transaction {
	return _trx.next
}

func (_trx *transaction) Push(t *transaction) {
	_trx.next = t
}
