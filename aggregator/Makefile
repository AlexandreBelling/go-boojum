build:
	cd c-boojum/build && cmake .. && make boojum

clear-build:
	mkdir -p ./c-boojum/build
	rm -rf ./c-boojum/build/*

clear-compiled:
	mkdir -p compiled
	rm -rf ./compiled/*

fetch-compiled:
	find . -name libboojum.a -exec cp {} ./compiled \;
	find . -name libff.a -exec cp {} ./compiled \;
	find . -name libsnark.a -exec cp {} ./compiled \;

clear-setup:
	mkdir -p setup
	rm -rf ./setup/*

run-test: clear-setup
	go test ; notify-send 'Done testing'

build-all: clear-build build clear-compiled fetch-compiled

all: build-all run-test
