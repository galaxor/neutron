all: build
start:
	go run neutron.go
build:
	go build -o neutron neutron.go
build-client:
	cd public && \
	sed -i 's/"angular"\: "1.4.9"/"angular": "~1.4.9"/g' bower.json && \
	sed -i "s/'nggettext_extract',/\/\/'nggettext_extract',/g" Gruntfile.js && \
	npm install && \
	node_modules/.bin/grunt ngconstant:dev build
clean-client-dist:
	rm -rf public/node_modules public/test
