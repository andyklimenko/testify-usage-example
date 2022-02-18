# testify-usage-example

![build](https://github.com/andyklimenko/testify-usage-example/actions/workflows/go.yml/badge.svg)

This is just an example of best testing practices I've come up with during my 5 yrs of professional Golang experience.
Feel free to report any issue or submit an improvement in PR if you find any

## Database
### local env
Use dockerized postresql for your convenience
```
docker run -p 5432:5432 -e POSTGRES_PASSWORD=secretpassword -d postgres
```

## Starting the webserver
At server's startup all db migrations will automatically be applied
```
export STORAGE_DRIVER=postgres
export STORAGE_DSN='user=postgres password=secretpassword dbname=postgres host=localhost port=5432 sslmode=disable'
export SERVER_ADDRESS=0.0.0.0:8080
export SERVER_NOTIFY_ADDRESS=0.0.0.0:8080
make run-server
```