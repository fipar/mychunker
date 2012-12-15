/*
chunk a mysql table and dump the chunks, optionally in parallel. 
Proof of concept. 
*/
package main

import (
	"flag"
	"fmt"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"os"
	"strconv"
	"sync"
)

var (
	verbose                                         bool
	host, user, password, port, schema, table, path string
	chunkSize, threads, running                     int
	db                                              *autorc.Conn
	lock                                            sync.Mutex
	wg                                              sync.WaitGroup
)

func main() {
	processArgs()
	initMySQL()
	cc, min, max := getChunkData()
	debug("Will chunk on " + cc + ", min = " + strconv.Itoa(min) + ", max = " + strconv.Itoa(max))
	dumpTable(cc, min, max)
	wg.Wait()
}

// Parses command line arguments
func processArgs() {
	flag.StringVar(&host, "host", "localhost", "The database host. Defaults to localhost")
	flag.StringVar(&user, "user", "root", "The database user. Defaults to root")
	flag.StringVar(&password, "password", "", "The database password. Defaults to empty string")
	flag.StringVar(&schema, "schema", "test", "The database schema. Defaults to test")
	flag.StringVar(&table, "table", "test", "The table to dump. Defaults to test")
	flag.StringVar(&port, "port", "3306", "The database port. Defaults to 3306")
	flag.StringVar(&path, "path", ".", "The path where the files will be created. Defaults to .")
	flag.IntVar(&chunkSize, "chunkSize", 1000, "The chunk size. Defaults to 1000")
	flag.IntVar(&threads, "threads", 4, "The number of threads. Defaults to 4")
	flag.BoolVar(&verbose, "verbose", false, "Be verbose")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

// so that I can do the same as panic() but without printing the stack trace to the user
func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}

// helper function to print debug messages
func debug(message string) {
	if verbose {
		fmt.Println("DEBUG: " + message)
	}
}

// dumps the table, one chunk at a time
func dumpTable(cc string, min int, max int) {
	running = 0
	for i := min; i < max; i += chunkSize {
		shouldGo := false
		lock.Lock()
		if running < threads {
			running++
			shouldGo = true
		}
		lock.Unlock()
		if shouldGo {
			go dumpChunk(cc, i, i+chunkSize-1, true)
		} else {
			dumpChunk(cc, i, i+chunkSize-1, false)
		}
	}
}

// dumps a specific chunk
// chunk will be delimited by cc (chunk column), values bewteen lower and upper
func dumpChunk(cc string, lower int, upper int, isThread bool) {
	var out *os.File
	dbcon := db
	defer func() {
		if isThread {
			lock.Lock()
			running--
			lock.Unlock()
			wg.Done()
		}
		out.Close()
		if err := recover(); err != nil {
			die(err.(error).Error())
		}
	}()
	if isThread {
		wg.Add(1)
		//dbcon = autorc.New("tcp", "", host+":"+port, user, password, schema)
	}
	out, _ = os.Create(path + "/" + schema + "." + table + "." + strconv.Itoa(lower) + "." + strconv.Itoa(upper) + ".csv")
	rows, _, _ := dbcon.Query("select * from " + schema + "." + table + " where " + cc + " between " + strconv.Itoa(lower) + " and " + strconv.Itoa(upper))
	for _, row := range rows {
		line := ""

		for idx, _ := range row {
			comma := ","
			if idx == len(row)-1 {
				comma = ""
			}
			line += row.Str(idx) + comma
		}

		out.WriteString(line + "\n")
	}
}

// gets the column that will be used for chunking, it's min and max values. 
// currently, this is the PK. if we make it pass the POC phase, it will be something better, akin to what pt-table-checksum uses
func getChunkData() (string, int, int) {
	defer func() {
		if err := recover(); err != nil {
			die(err.(error).Error())
		}
	}()
	rows, _, _ := db.Query("select column_name from information_schema.columns where table_schema='" + schema + "' and table_name='" + table + "' and column_key='PRI'")
	if rows == nil {
		die("Could not find a column suitable for chunking")
	}
	cc := rows[0].Str(0)
	rows, _, _ = db.Query("select min(" + cc + "), max(" + cc + ") from " + schema + "." + table)
	return cc, rows[0].Int(0), rows[0].Int(1)
}

// connects to MySQL
func initMySQL() {
	var err error
	debug("connecting to MySQL")
	debug("Will connect to tcp, " + host + ":" + port + ", " + user + ", " + password + ", " + schema)
	db = autorc.New("tcp", "", host+":"+port, user, password, schema)
	_, _, err = db.Query("select 1")
	if err != nil {
		die("An error occurred while connecting to MySQL: " + err.Error())
	}
}
