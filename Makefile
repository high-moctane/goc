build: goc

dockerbuild:
	make clean
	docker run --rm -v $(CURDIR):/goc -w /goc goc make build

clean:
	rm -rf goc *.o *~ tmp*

dockerrun:
	docker run --rm -it -v $(CURDIR):/goc -w /goc goc

dockertest:
	make dockerbuild
	docker run --rm -v $(CURDIR):/goc -w /goc goc make test

docker:
	docker build -t goc .

test: goc
	./test.sh

goc: main.go
	go build

.PHONY: build dockerbuild clean dockerrun dockertest docker test