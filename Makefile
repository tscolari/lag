all:
	GOOS=linux go build -o lag-linux main.go
	GOOS=darwin go build -o lag-darwin main.go
