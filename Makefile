all: build
start:
	go run neutron.go
build:
	go build -o neutron neutron.go
build-client:
	cd public && \
	npm install --ignore-scripts && \
	bower install -F && \
	sed -i 's/https:\/\/dev\.protonmail\.com//g' src/app/config.js && \
	node_modules/.bin/grunt build
