package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson/http"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	dsn := flag.String("dsn", ":memory:", "")
	port := flag.Int("port", 8080, "")
	host := flag.String("host", "localhost", "")

	flag.Parse()

	db, err := database.NewDBWithDriver("spatialite", *dsn)

	if err != nil {
		log.Fatal(err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	nearby_handler, err := http.NearbyHandler(db)

	if err != nil {
		log.Fatal(err)
	}

	pip_handler, err := http.PointInPolygonHandler(db)

	if err != nil {
		log.Fatal(err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/ping", ping_handler)
	mux.Handle("/nearby", nearby_handler)
	mux.Handle("/pip", pip_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening for requests on %s\n", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
