package gdbc

import (
	"database/sql"
	"database/sql/driver"
	"sync"
)

const GdbcDriverPrefix = "gdbc$"

type DataSourceNameAdapter interface {
	GetDataSourceName(dataSource DataSource) (string, error)
}

var (
	dsnAdapterMu sync.RWMutex
	dsnAdapters  = make(map[string]DataSourceNameAdapter)
)

func Register(name string, driver driver.Driver, dsnAdapter DataSourceNameAdapter) {
	dsnAdapterMu.Lock()
	defer dsnAdapterMu.Unlock()
	sql.Register(GdbcDriverPrefix+name, driver)

	if dsnAdapter == nil {
		panic("sql: DSN adapter is nil")
	}
	dsnAdapters[GdbcDriverPrefix+name] = dsnAdapter
}
