pwd
go mod vendor
go mod tidy
go build -v -o ./dummy/bin/dummy ./dummy/replica/
go build -v -o ./torture/bin/torture ./torture/torture/
