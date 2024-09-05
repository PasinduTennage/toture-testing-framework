pwd
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build -v -o ./consenbench/bin/bench ./consenbench/
