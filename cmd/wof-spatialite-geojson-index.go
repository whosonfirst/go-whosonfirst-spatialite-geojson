package main

import (
	"flag"
	"fmt"
	wof_index "github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson/index"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	sql_index "github.com/whosonfirst/go-whosonfirst-sqlite/index"
	"io"
	"os"
	"runtime"
	"strings"
)

func main() {

	valid_modes := strings.Join(wof_index.Modes(), ",")
	desc_modes := fmt.Sprintf("The mode to use importing data. Valid modes are: %s.", valid_modes)

	mode := flag.String("mode", "files", desc_modes)

	dsn := flag.String("dsn", ":memory:", "")

	all := flag.Bool("all", false, "Index all tables")
	geojson := flag.Bool("geojson", true, "Index the 'geojson' table")
	live_hard := flag.Bool("live-hard-die-fast", true, "Enable various performance-related pragmas at the expense of possible (unlikely) database corruption")
	timings := flag.Bool("timings", false, "Display timings during and after indexing")
	var procs = flag.Int("processes", (runtime.NumCPU() * 2), "The number of concurrent processes to index data with")

	is_wof := flag.Bool("is-wof", true, "...")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	logger := log.SimpleWOFLogger()

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, "status")

	db, err := database.NewDBWithDriver("spatialite", *dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	if *live_hard {

		err = db.LiveHardDieFast()

		if err != nil {
			logger.Fatal("Unable to live hard and die fast so just dying fast instead, because %s", err)
		}
	}

	to_index := make([]sqlite.Table, 0)

	if *geojson || *all {

		gt, err := tables.NewGeoJSONTableWithDatabase(db)

		if err != nil {
			logger.Fatal("failed to create 'geojson' table because %s", err)
		}

		to_index = append(to_index, gt)
	}

	if len(to_index) == 0 {
		logger.Fatal("You forgot to specify which (any) tables to index")
	}

	var idx *sql_index.SQLiteIndexer

	if *is_wof {

		i, err := index.NewSpatialiteWOFIndexer(db, to_index)

		if err != nil {
			logger.Fatal("failed to create sqlite indexer because %s", err)
		}

		idx = i

	} else {

		i, err := index.NewSpatialiteGeoJSONIndexer(db, to_index)

		if err != nil {
			logger.Fatal("failed to create sqlite indexer because %s", err)
		}

		idx = i
	}

	idx.Timings = *timings
	idx.Logger = logger

	err = idx.IndexPaths(*mode, flag.Args())

	if err != nil {
		logger.Fatal("Failed to index paths in %s mode because: %s", *mode, err)
	}

	os.Exit(0)
}
