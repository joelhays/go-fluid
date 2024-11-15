all: clean build
clean:
		go clean ./...
		rm -rf ./bin
		rm -f ./main
test:
		go test -v ./...
build:
		go build -o ./bin/app -v
run: build
		./bin/app $(ARGS)
