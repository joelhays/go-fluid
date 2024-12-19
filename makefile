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
		./bin/app
runprof: build
		./bin/app --profile
pprof:
		go tool pprof -http=:8080 go-fluid.prof
