package store

type IStore interface {
	Begin() error
	Rollback() error
	Commit() error
	Count(value string) int
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}
