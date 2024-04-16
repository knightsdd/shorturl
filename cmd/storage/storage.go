package storage

type MyStorage map[string]string

var storage MyStorage = make(MyStorage, 10)

func GetStorage() MyStorage {
	return storage
}
