DESTDIR ?=
PREFIX ?= /usr/local
OUTPUT_DIR ?= out
bin = wakeup
outfile = $(OUTPUT_DIR)/$(bin).linux-$(shell uname -m)

build:
	go build -o $(outfile) ./main.go
	sudo setcap cap_net_raw+ep $(outfile)

run:
	$(outfile)

dev: build run


clean:
	rm -rf out

deps:
	go get ./...