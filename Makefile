all: 
	go build mychunker.go

run: all
	./mychunker -port 5527 -user msandbox -password msandbox -verbose -schema test -table test

clean: 
	rm -f mychunker test.* 
