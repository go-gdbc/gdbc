package gdbc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDataSourceUrlWithWrongFormat(t *testing.T) {
	dataSourceUrl := "wrong-format"
	dataSource, err := parse(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "url format is wrong : "+dataSourceUrl, err.Error())

	dataSourceUrl = "gdbc:driver-name:"
	dataSource, err = parse(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "url format is wrong : "+dataSourceUrl, err.Error())
}

func TestParseDataSourceUrlNonStartingWithGdbc(t *testing.T) {
	dataSourceUrl := "test:driver:db"
	dataSource, err := parse(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "url must start with 'gdbc'", err.Error())
}

func TestParseDataSourceUrlWithEmptyDriverName(t *testing.T) {
	dataSourceUrl := "gdbc::db"
	dataSource, err := parse(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "driver name must not be empty", err.Error())
}

func TestParseDataSourceUrlWithOpaque(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name:db"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "db", dataSource.url.Opaque)
}

func TestParseDataSourceUrlWithOpaqueAndArgs(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name:db?arg1=value1&arg2=value2"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "db", dataSource.url.Opaque)
	assert.Equal(t, "value1", dataSource.url.Query().Get("arg1"))
	assert.Equal(t, "value2", dataSource.url.Query().Get("arg2"))
	assert.Equal(t, "", dataSource.url.Query().Get("arg3"))
}

func TestParseDataSourceUrlWithHost(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://localhost"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "localhost", dataSource.url.Host)
}

func TestParseDataSourceUrlWithHostAndArgs(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://localhost?arg1=value1&arg2=value2"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "localhost", dataSource.url.Host)
	assert.Equal(t, "value1", dataSource.url.Query().Get("arg1"))
	assert.Equal(t, "value2", dataSource.url.Query().Get("arg2"))
	assert.Equal(t, "", dataSource.url.Query().Get("arg3"))
}

func TestParseDataSourceUrlWithHostAndPort(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://localhost:5432"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "localhost:5432", dataSource.url.Host)
	assert.Equal(t, "localhost", dataSource.url.Hostname())
	assert.Equal(t, "5432", dataSource.url.Port())
}

func TestParseDataSourceUrlWithHostAndPortAndArgs(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://localhost:5432?arg1=value1&arg2=value2"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "localhost:5432", dataSource.url.Host)
	assert.Equal(t, "localhost", dataSource.url.Hostname())
	assert.Equal(t, "5432", dataSource.url.Port())
	assert.Equal(t, "value1", dataSource.url.Query().Get("arg1"))
	assert.Equal(t, "value2", dataSource.url.Query().Get("arg2"))
	assert.Equal(t, "", dataSource.url.Query().Get("arg3"))
}

func TestParseDataSourceUrlWithHostAndUserInfo(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://username:password@localhost:5432"
	dataSource, err := parse(dataSourceUrl)
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.driverName)

	assert.NotNil(t, dataSource.url)
	assert.Equal(t, "localhost:5432", dataSource.url.Host)
	assert.Equal(t, "localhost", dataSource.url.Hostname())
	assert.Equal(t, "5432", dataSource.url.Port())

	assert.NotNil(t, dataSource.url.User)
	assert.Equal(t, "username", dataSource.url.User.Username())
	password, _ := dataSource.url.User.Password()
	assert.Equal(t, "password", password)
}

func TestParseDataSourceUrlWithWrongHostFormat(t *testing.T) {
	dataSourceUrl := "gdbc:driver-name://localhost:port:wtf"
	dataSource, err := parse(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "parse gdbc://localhost:port:wtf: invalid port \":wtf\" after host", err.Error())
}

func TestGetDataSource(t *testing.T) {
	dataSource, err := GetDataSource("gdbc:driver-name://localhost:5432")
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.GetDriverName())
	assert.NotNil(t, dataSource.GetURL())
	assert.Equal(t, "localhost", dataSource.url.Hostname())
	assert.Equal(t, "5432", dataSource.url.Port())
	assert.Empty(t, dataSource.GetUsername())
	assert.Empty(t, dataSource.GetPassword())
}

func TestGetDataSourceWithUsernameAndPassword(t *testing.T) {
	dataSource, err := GetDataSource("gdbc:driver-name://localhost:5432", Username("username"), Password("password"))
	assert.NotNil(t, dataSource)
	assert.Nil(t, err)

	assert.Equal(t, "driver-name", dataSource.GetDriverName())
	assert.NotNil(t, dataSource.GetURL())
	assert.Equal(t, "localhost", dataSource.url.Hostname())
	assert.Equal(t, "5432", dataSource.url.Port())
	assert.Equal(t, "username", dataSource.GetUsername())
	assert.Equal(t, "password", dataSource.GetPassword())
}

func TestGetDataSourceUrlWithWrongFormat(t *testing.T) {
	dataSourceUrl := "wrong-format"
	dataSource, err := GetDataSource(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "url format is wrong : "+dataSourceUrl, err.Error())

	dataSourceUrl = "gdbc:driver-name:"
	dataSource, err = GetDataSource(dataSourceUrl)
	assert.Nil(t, dataSource)
	assert.NotNil(t, err)
	assert.Equal(t, "url format is wrong : "+dataSourceUrl, err.Error())
}
