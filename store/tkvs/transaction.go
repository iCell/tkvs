package tkvs

type transaction struct {
	Kvs  map[string]string
	Next *transaction
}

func newTransaction() *transaction {
	return &transaction{Kvs: make(map[string]string)}
}

func (_trx *transaction) Set(k, v string) {
	_trx.Kvs[k] = v
}

func (_trx *transaction) Get(k string) (string, bool) {
	v, exist := _trx.Kvs[k]
	return v, exist
}

func (_trx *transaction) Delete(k string) {
	delete(_trx.Kvs, k)
}

func (_trx *transaction) Clear() {
	for key := range _trx.Kvs {
		delete(_trx.Kvs, key)
	}
}
