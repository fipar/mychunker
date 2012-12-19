all: 
	go build mychunker.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build mychunker.go

# my run and run_single targets assume mysql will listen on 5527 since I'm using a sandbox. 
run: all
	./mychunker -port 5527 -user msandbox -password msandbox -verbose -schema test -table test

run_single: all
	./mychunker -port 5527 -user msandbox -password msandbox -verbose -schema test -table test -threads 1

# assumes the mysql cli is configured (proper ~/.my.cnf) to access the same server run and run_single will
init_table: 
	./init_table.sh
clean: 
	rm -f mychunker test.* 

clean-data:
	rm -f test.*
