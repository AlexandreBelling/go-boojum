build:
	cd c-boojum/build && cmake .. && make boojum

clear-build:
	mkdir -p ./c-boojum/build
	rm -rf ./c-boojum/build/*

clear-compiled:
	mkdir -p compiled
	rm -rf ./compiled/*

fetch-compiled:
	find . -name \*.a -exec cp {} ./compiled \;
	find . -name \*.o -exec cp {} ./compiled \;
	find . -name \*.so -exec cp {} ./compiled \;

clear-setup:
	mkdir -p setup
	rm -rf ./setup/*

run-test: clear-setup
	go test ; notify-send 'Done testing'

build-all: clear-build build clear-compiled fetch-compiled

all: clear-build build clear-compiled fetch-compiled run-test
