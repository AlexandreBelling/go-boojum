all:
	git submodule update --init --recursive
	cd aggregator && make build-all