all: build
start:
	go run neutron.go
build:
	go build -o neutron neutron.go
