package gdbc

import (
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	Register("testDriver1", &testDriver1{}, &testDriver1DataSourceNameAdapter{})
	Register("testDriver2", &testDriver2{}, &testDriver2DataSourceNameAdapter{})
}

type testConnection struct {
}

func (conn testConnection) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (conn testConnection) Close() error {
	return nil
}

func (conn testConnection) Begin() (driver.Tx, error) {
	return nil, nil
}

type testDriver1 struct {
}

func (driver testDriver1) Open(name string) (driver.Conn, error) {
	return &testConnection{}, nil
}

type testDriver1DataSourceNameAdapter struct {
}

func (adapter testDriver1DataSourceNameAdapter) GetDataSourceName(dataSource DataSource) (string, error) {
	return "test", nil
}

type testDriver2 struct {
}

func (driver testDriver2) Open(name string) (driver.Conn, error) {
	return &testConnection{}, nil
}

type testDriver2DataSourceNameAdapter struct {
}

func (adapter testDriver2DataSourceNameAdapter) GetDataSourceName(dataSource DataSource) (string, error) {
	return "", errors.New("test error")
}

func TestRegisterWithNilAdapter(t *testing.T) {
	assert.Panics(t, func() {
		Register("testDriver", &testDriver1{}, nil)
	})
}
