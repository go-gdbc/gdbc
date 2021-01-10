package gdbc

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DataSourceOption func(dataSource *DataSource)

const Scheme = "gdbc"

type DataSource struct {
	driverName string
	url        *url.URL
	username   string
	password   string
}

func (dataSource *DataSource) GetDriverName() string {
	return dataSource.driverName
}

func (dataSource *DataSource) GetURL() *url.URL {
	return dataSource.url
}

func (dataSource *DataSource) GetUsername() string {
	return dataSource.username
}

func (dataSource *DataSource) GetPassword() string {
	return dataSource.password
}

func (dataSource *DataSource) GetConnection() (*sql.DB, error) {
	driversMu.RLock()
	driver, ok := drivers[dataSource.driverName]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", dataSource.driverName)
	}

	dataSourceName, err := driver.GetDataSourceName(dataSource)
	if err != nil {
		return nil, err
	}

	return sql.Open(dataSource.driverName, dataSourceName)
}

func GetDataSource(uri string, options ...DataSourceOption) (*DataSource, error) {
	dataSource, err := parse(uri)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(dataSource)
	}
	return dataSource, nil
}

func Username(username string) DataSourceOption {
	return func(dataSource *DataSource) {
		dataSource.username = username
	}
}

func Password(password string) DataSourceOption {
	return func(dataSource *DataSource) {
		dataSource.password = password
	}
}

func parse(uri string) (*DataSource, error) {
	src := strings.Split(uri, ":")
	if len(src) < 3 {
		return nil, errors.New("uri format is wrong : " + uri)
	}

	scheme := src[0]
	if Scheme != scheme {
		return nil, errors.New("uri must start with 'gdbc'")
	}

	driverName := src[1]
	if len(driverName) == 0 {
		return nil, errors.New("driver name must not be empty")
	}

	dataSource := &DataSource{
		driverName: driverName,
	}
	rest := strings.Join(append(src[:1], src[2:]...), ":")
	if len(rest) != 0 {
		parsedUrl, err := url.ParseRequestURI(rest)
		if err != nil {
			return nil, err
		}
		dataSource.url = parsedUrl
		return dataSource, nil
	}
	return nil, errors.New("uri format is not wrong : " + uri)
}
