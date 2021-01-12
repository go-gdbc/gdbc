# GDBC - Go Database Connectivity
Data Source Name, known as DSN does not have any standard format. 

With GDBC, a database is represented by a URL (Uniform Resource Locator).
URL takes one of the following forms:
```
gdbc:driver-name:database?arg1=value1&arg2=value...
gdbc:driver-name://localhost/database?arg1=value1&arg2=value...
gdbc:driver-name://localhost:5432/database?arg1=value1&arg2=value...
gdbc:driver-name://username:password@localhost:5432/database?arg1=value1&arg2=value...
gdbc:driver-name:file:h2?arg1=value1&arg2=value...
```

## How to use GDBC?

```go
dataSource, err := GetDataSource("gdbc:driver-name://localhost:5432/test-db", Username("username"), Password("password"))
if err != nil {
    panic(err)
}

var connection *sql.DB
connection, err = dataSource.GetConnection()
if err != nil {
    panic(err)
}
...

```

## Adapter
```go
type DataSourceNameAdapter interface {
	GetDataSourceName(dataSource *DataSource) (string, error)
}
```

## License
GDBC is released under version 2.0 of the Apache License