pwd
go mod vendor
go mod tidy
go build -v -o ./dummy/bin/dummy ./dummy/replica/
go build -v -o ./torture/bin/torture ./torture/torture/
go build -v -o ./test/bin/test ./test/
