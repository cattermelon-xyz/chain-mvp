build:
	go build -o ./bin/htg-client build/client/main.go
	go build -o ./bin/htg build/server/main.go
	# go build -o ./bin/demo build/main.go
serve:
	go build -o ./bin/htg build/server/main.go
	./bin/htg start
test:
	go test -v ./types/...
