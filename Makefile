all: build

build:
	go build .

clean:
	rm -f openshift-tests-endpoint