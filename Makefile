all: build
start:
	go run neutron.go
build:
	go build -o neutron neutron.go
build-client:
	cd public && npm install && sed -i 's/https:\/\/dev\.protonmail\.com//g' public/src/app/config.js && grunt build
