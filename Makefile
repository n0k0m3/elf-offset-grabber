BINARY_NAME=grabber
ATELIER_SRC=main.go

build:
	go build -o ./bin/$(BINARY_NAME).exe
	copy conf.toml bin\conf.toml

build-linux:
	go build -o ./bin/$(BINARY_NAME)
	cp ./conf.toml ./bin/conf.toml