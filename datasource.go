package gdbc

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DataSourceOption func(dataSource DataSource)

const Scheme = "gdbc"

type DataSource interface {
	GetDriverName() string
	GetURL() *url.URL
	GetConnection() (*sql.DB, error)
	GetUsername() string
	SetUsername(username string)
	GetPassword() string
	SetPassword(password string)
}

type SimpleDataSource struct {
	driverName string
	url        *url.URL
	username   string
	password   string
}

func (dataSource *SimpleDataSource) GetDriverName() string {
	return dataSource.driverName
}

func (dataSource *SimpleDataSource) GetURL() *url.URL {
	return dataSource.url
}

func (dataSource *SimpleDataSource) GetUsername() string {
	return dataSource.username
}

func (dataSource *SimpleDataSource) SetUsername(username string) {
	dataSource.username = username
}

func (dataSource *SimpleDataSource) GetPassword() string {
	return dataSource.password
}

func (dataSource *SimpleDataSource) SetPassword(password string) {
	dataSource.password = password
}

func (dataSource *SimpleDataSource) GetConnection() (*sql.DB, error) {
	dsnAdapter := GetDataSourceNameAdapter(dataSource.driverName)
	if dsnAdapter == nil {
		return nil, fmt.Errorf("sql: driver does not exist : %s", dataSource.driverName)
	}

	dataSourceName, err := dsnAdapter.GetDataSourceName(dataSource)
	if err != nil {
		return nil, err
	}

	driverName := driverNameMap[dataSource.driverName]
	return sql.Open(driverName, dataSourceName)
}

func GetDataSource(url string, options ...DataSourceOption) (DataSource, error) {
	dataSource, err := parse(url)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(dataSource)
	}
	return dataSource, nil
}

func Username(username string) DataSourceOption {
	return func(dataSource DataSource) {
		dataSource.SetUsername(username)
	}
}

func Password(password string) DataSourceOption {
	return func(dataSource DataSource) {
		dataSource.SetPassword(password)
	}
}

func parse(dataSourceUrl string) (DataSource, error) {
	src := strings.Split(dataSourceUrl, ":")
	if len(src) < 3 {
		return nil, errors.New("url format is wrong : " + dataSourceUrl)
	}

	scheme := src[0]
	if Scheme != scheme {
		return nil, errors.New("url must start with 'gdbc'")
	}

	driverName := src[1]
	if len(driverName) == 0 {
		return nil, errors.New("driver name must not be empty")
	}

	dataSource := &SimpleDataSource{
		driverName: driverName,
	}
	rest := strings.Join(append(src[:1], src[2:]...), ":")
	if rest != Scheme+":" {
		parsedUrl, err := url.ParseRequestURI(rest)
		if err != nil {
			return nil, err
		}
		dataSource.url = parsedUrl
		return dataSource, nil
	}
	return nil, errors.New("url format is wrong : " + dataSourceUrl)
}
