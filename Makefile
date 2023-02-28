build:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/SCUtil.exe main.go

optimised-build:
	go build -ldflags "-s -w"