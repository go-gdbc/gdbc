# GDBC - Go Database Connectivity
Because Data Source Name, known as DSN does not have any standard format, the driver libraries have their 
driver-specific DSN. Sometimes you might get confused about how to specify your DSN for the database you 
want to connect to. In order to solve this issue, GDBC provides a connection format to represent the database, 
and an abstract layer for database connections. 

Note that drivers have to be registered by using **gdbc.Register** instead **sql.Register**.
Otherwise, you cannot connect to the database by using gdbc.

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
It's pretty easy to connect the database. You can connect to the database as shown below.

```go
dataSource, err := gdbc.GetDataSource("gdbc:driver-name://localhost:5432/test-db", Username("username"), Password("password"))
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

## Register your existing driver
If you have already an existing driver written and want to support GDBC, you have to implement the following 
**DataSourceNameAdapter** interface converting GDBC format to your specific DSN format, and register your
driver by using **gdbc.Register** function.

```go
type DataSourceNameAdapter interface {
	GetDataSourceName(dataSource DataSource) (string, error)
}
```

## License
GDBC is released under version 2.0 of the Apache License