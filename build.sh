pwd
sudo apt install -y protobuf-compiler
/usr/local/go/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
/usr/local/go/bin/go get github.com/shirou/gopsutil/v3/...
/usr/local/go/bin/go get fyne.io/fyne/v2
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
protoc --go_out=. --go_opt=paths=source_relative consenbench/common/message.proto
/usr/local/go/bin/go build -v -o ./consenbench/bin/bench ./consenbench/