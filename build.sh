pwd
go mod vendor
go mod tidy
go build -v -o ./toture/bin/toture ./toture/
go build -v -o ./dummy/bin/dummy ./dummy/