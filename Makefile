build:
	GOOS=darwin GOARCH=amd64 go build -o hex

run:
	./hex

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o hex && ./hex