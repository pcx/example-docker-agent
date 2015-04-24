linux:
	GOOS=linux GOARCH=amd64 go build  -o bin/stacktape-agent_linux-amd64 enf/enf.go
darwin:
	GOOS=darwin GOARCH=amd64 go build  -o bin/stacktape-agent_darwin-amd64 enf/enf.go
all:
	make clean
	make linux
	make darwin
clean:
	rm -r bin/