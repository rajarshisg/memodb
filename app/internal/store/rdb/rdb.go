package rdb

type KVValue struct {
	Value string;
	ExpireAt uint64;
}
type RDBDatabase struct {
	DatabaseNumber int;
	HashTableSize int;
	ExpiryHashTableSize int;
	KVMap map[string]KVValue;
}
type RDBType struct {
	Version string;
	Databases []RDBDatabase;
}