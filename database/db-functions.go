package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DBVertexResult struct {
	vertexID     uint64
	vertexIDHash uint64
	neighbors    []uint64
}

var db *sql.DB
var dbName = "bagelDB_new"
var tableName = "adjList"
var server = "bagel.database.windows.net"
var port = 1433
var user = "user"
var password = "Distributedgraph!"
var database = "bagel_2.0"

func main() {
	start := time.Now()
	n, err := getVerticesModulo(1, 3)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Found vertex: %v in %s\n", n, elapsed)
}

func connectToDb() (*sql.DB, error) {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return nil, err
	}
	return db, nil
}

func getVertex(id int) (*DBVertexResult, error) {
	connectToDb()
	if db == nil {
		fmt.Println("Not connected to Database yet")
		panic("aaa")
	}
	rows, err := db.Query("SELECT * FROM " + tableName + " WHERE srcVertex = " + strconv.Itoa(id) + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Query for a value based on a single row.
	var searchID uint64
	var hash string
	var neighbors string
	qs := "SELECT * FROM " + tableName + " WHERE srcVertex = " + strconv.Itoa(id) + ";"
	if err := db.QueryRow(qs).Scan(&searchID, &hash, &neighbors); err != nil {
		if err == sql.ErrNoRows {
			return &DBVertexResult{}, fmt.Errorf("%d: unknown ID", id)
		}
		return &DBVertexResult{}, fmt.Errorf("some kind of error :| %d", id)
	}
	hashNum, err := strconv.ParseUint(hash, 10, 64)
	if err != nil {
		panic("parsing hash to Uint64 failed")
	}
	v := DBVertexResult{vertexID: searchID, vertexIDHash: hashNum, neighbors: stringToArray(neighbors, ".")}
	return &v, nil
}

func getVerticesModulo(workerId uint32, numWorkers uint8) ([]DBVertexResult, error) {
	connectToDb()
	if db == nil {
		fmt.Println("Not connected to Database yet")
		panic("aaa")
	}
	result, err := db.Query(
		"SELECT * from " + tableName + " where srcVertex % " + strconv.Itoa(int(numWorkers)) + " = " + strconv.Itoa(int(workerId)) + ";")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic("query went wrong")
	}
	var searchID uint64
	var hash string
	var neighbors string
	var vertices = []DBVertexResult{}
	for result.Next() {
		err := result.Scan(&searchID, &hash, &neighbors)
		if err != nil {
			panic("scan went wrong")
		}
		hashNum, err := strconv.ParseUint(hash, 10, 64)
		if err != nil {
			panic("parsing hash to Uint64 failed")
		}
		v := DBVertexResult{vertexID: searchID, vertexIDHash: hashNum, neighbors: stringToArray(neighbors, ".")}
		vertices = append(vertices, v)
	}

	return vertices, nil
}

func stringToArray(a string, delim string) []uint64 {
	neighbors := strings.Split(a, delim)
	neighborSlice := []uint64{}
	if len(strings.TrimSpace(a)) == 0 {
		return neighborSlice
	}
	for _, v := range neighbors {
		neighborID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic("parsing hash to Uint64 failed")
		}
		neighborSlice = append(neighborSlice, neighborID)
	}
	return neighborSlice
}
