DESTDIR ?=
PREFIX ?= /usr/local
OUTPUT_DIR ?= out
bin = wakeup

build:
	go build -o $(OUTPUT_DIR)/$(bin) ./main.go
	sudo setcap cap_net_raw+ep $(OUTPUT_DIR)/$(bin)

run:
	$(OUTPUT_DIR)/$(bin)

dev: build run


clean:
	rm -rf out

deps:
	go get ./...