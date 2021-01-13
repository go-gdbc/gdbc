package gdbc

import (
	"sync"
)

type DataSourceNameAdapter interface {
	GetDataSourceName(dataSource DataSource) (string, error)
}

var (
	dsnAdapterMu  sync.RWMutex
	dsnAdapters   = make(map[string]DataSourceNameAdapter)
	driverNameMap = make(map[string]string)
)

func Register(driverName string, driverAliasName string, dsnAdapter DataSourceNameAdapter) {
	dsnAdapterMu.Lock()
	defer dsnAdapterMu.Unlock()

	if driverName == "" || driverAliasName == "" {
		panic("sql: Driver and driver alias name cannot be empty")
	}

	if dsnAdapter == nil {
		panic("sql: DSN adapter is nil")
	}

	_, ok := driverNameMap[driverAliasName]
	if ok {
		panic("sql: already registered an adapter with the same alias name : " + driverAliasName)
	}
	driverNameMap[driverAliasName] = driverName

	_, ok = dsnAdapters[driverName]
	if ok {
		panic("sql: you already registered a DSN adapter for driver : " + driverName)
	}
	dsnAdapters[driverName] = dsnAdapter
}

func GetDataSourceNameAdapter(driverAliasName string) DataSourceNameAdapter {
	dsnAdapterMu.RLock()
	driverName, ok := driverNameMap[driverAliasName]
	dsnAdapterMu.RUnlock()
	if !ok {
		return nil
	}
	return dsnAdapters[driverName]
}

func GetDriverName(driverAliasName string) (string, bool) {
	dsnAdapterMu.RLock()
	driverName, ok := driverNameMap[driverAliasName]
	dsnAdapterMu.RUnlock()
	if !ok {
		return "", false
	}
	return driverName, true
}
