SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./ele_document main.go
pause