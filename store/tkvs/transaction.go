package tkvs

type storedValue struct {
	Val   string
	Valid bool
}

type transaction struct {
	Kvs  map[string]*storedValue
	Next *transaction
}

func newTransaction() *transaction {
	return &transaction{Kvs: make(map[string]*storedValue)}
}

func (_trx *transaction) Set(k string, v *storedValue) {
	_trx.Kvs[k] = v
}

func (_trx *transaction) Get(k string) (*storedValue, bool) {
	v, exist := _trx.Kvs[k]
	return v, exist
}

func (_trx *transaction) Delete(k string) {
	_, exist := _trx.Get(k)
	if exist {
		_trx.Kvs[k].Valid = false
		return
	}
	_trx.Kvs[k] = &storedValue{
		Val:   "",
		Valid: false,
	}
}
