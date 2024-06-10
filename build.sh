pwd
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build -v -o ./dummy/bin/dummy ./dummy/replica/
/usr/local/go/bin/go build -v -o ./torture/bin/torture ./torture/torture/
