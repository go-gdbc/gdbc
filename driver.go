package gdbc

import (
	"database/sql"
	"database/sql/driver"
	"sync"
)

type DataSourceNameAdapter interface {
	GetDataSourceName(dataSource *DataSource) (string, error)
}

type Driver interface {
	driver.Driver
	DataSourceNameAdapter
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	sql.Register(name, driver)
	drivers[name] = driver
}
