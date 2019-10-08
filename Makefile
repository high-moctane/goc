build: goc

clean:
	rm -rf goc *.o *~ tmp*

dockerrun:
	docker run --rm -it -v $(CURDIR):/goc -w /goc goc

dockertest:
	docker run --rm -it -v $(CURDIR):/goc -w /goc goc make test

docker:
	docker build -t goc .

test: goc
	./test.sh

goc: main.go
	go build

.PHONY: build dockerrun dockertest docker test