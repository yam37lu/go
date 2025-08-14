SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GIN_MODE=release
go build  -o ../target/bin/template ../src/main.go
