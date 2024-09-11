pwd
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/shirou/gopsutil/v3/...
go get fyne.io/fyne/v2
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
protoc --go_out=. --go_opt=paths=source_relative consenbench/common/message.proto
/usr/local/go/bin/go build -v -o ./consenbench/bin/bench ./consenbench/
