all:
	@go get ./src
	@go build -o hook ./src

run:
	@go run ./src
