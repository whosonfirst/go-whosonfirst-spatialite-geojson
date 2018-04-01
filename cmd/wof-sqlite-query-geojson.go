package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-sqlite-geojson/query"
	"github.com/whosonfirst/go-whosonfirst-sqlite-geojson/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"io"
	"os"
)

func main() {

	dsn := flag.String("dsn", ":memory:", "")
	limit := flag.Int("limit", 10, "")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, "status")

	db, err := database.NewDBWithDriver("spatialite", *dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	conn, err := db.Conn()

	if err != nil {
		logger.Fatal("unable to create database connection because %s", err)
	}

	t, err := tables.NewGeoJSONTable()

	if err != nil {
		logger.Fatal("unable to create table because %s", err)
	}

	args := make([]interface{}, 0)

	q := fmt.Sprintf("SELECT id, properties, AsGeoJSON(geometry) FROM %s", t.Name())

	if *limit > 0 {
		q = fmt.Sprintf("%s LIMIT %d", q, *limit)
	}

	fc, err := query.QueryToFeatureCollection(conn, q, args...)

	if err != nil {
		logger.Fatal("failed to query (%s) because %s", q, err)
	}

	js, err := json.Marshal(fc)

	if err != nil {
		logger.Fatal("failed to seraliaze feature collection %s", err)
	}

	fmt.Println(string(js))
	os.Exit(0)
}
