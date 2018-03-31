# go-whosonfirst-sqlite-geojson

Go package for indexing GeoJSON features in a Spatialite (SQLite) database.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tables

### geojson

```
CREATE TABLE geojson (
	id INTEGER NOT NULL PRIMARY KEY,
	properties JSON
);

SELECT InitSpatialMetaData();
SELECT AddGeometryColumn('geojson', 'geom', 4326, 'GEOMETRY', 'XY');
SELECT CreateSpatialIndex('geojson', 'geom');

```

In order to index geometries you will need to have the [Spatialite extension](https://www.gaia-gis.it/fossil/libspatialite/index) already installed. Installation details are outside the scope of this document.

This package also assumes you if you are accessing the resulting database using a `sqlite3` binary or language-specific SQLite3 library (other than this package) that the [JSON1 extension](https://www.sqlite.org/json1.html) is available and loaded.

## Custom tables

Sure. You just need to write a per-table package that implements the `Table` interface as described in [go-whosonfirst-sqlite](https://github.com/whosonfirst/go-whosonfirst-sqlite#custom-tables).

## Tools

### wof-sqlite-index-geojson

```
./bin/wof-sqlite-index-geojson -h
Usage of ./bin/wof-sqlite-index-geojson:
  -all
	Index all tables (except the 'search' and 'geometries' tables which you need to specify explicitly)
  -dsn string
        (default ":memory:")
  -geojson
	Index the 'geojson' table (default true)
  -live-hard-die-fast
	Enable various performance-related pragmas at the expense of possible (unlikely) database corruption (default true)
  -mode string
    	The mode to use importing data. Valid modes are: directory,feature,feature-collection,files,geojson-ls,meta,path,repo,sqlite. (default "files")
  -processes int
    	     The number of concurrent processes to index data with (default 8)
  -timings
	Display timings during and after indexing
```

_Please finish writing me..._

## Querying the data

There isn't a general purpose query tool as of this writing. Every GeoJSON
record's `properties` dictionary is stored as a SQLite `JSON` database column
and it is assumed you will query them using the `json_extract` function.

Let's say that you've indexed the
[whosonfirst-data-constituency-ca](https://github.com/whosonfirst-data/whosonfirst-data-constituency-ca)
repo like this:

```
./bin/wof-sqlite-index-geojson -dsn test.db -mode repo /usr/local/data/whosonfirst-data-constituency-ca/
```

You might then query the data (in a Go program or anything else that can talk to SQLite) like this:

```
package main

import (
        "github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"log"
)

func main (){

	db, _ := database.NewDBWithDriver("spatialite", "test.db")
        defer db.Close()

        conn, _ := db.Conn()

	sql := "SELECT json_extract(properties, '$.ebc:ed_abbrev') AS ed FROM geojson WHERE json_extract(properties, '$.ebc:ed_abbrev') = 'ABM'"
        row := conn.QueryRow(sql)

	var eb string
	row.Scan(&eb)

        log.Println(eb)
}
```

## Modes

This package can index any input source supported by the
[go-whosonfirst-index](https://github.com/whosonfirst/go-whosonfirst-index#modes)
package.

## Performance

Unknown. On the face of it, it seems certain that having to call `json_extract`
_all the time_ wouldn't be terribly performant. Think of it as a cheap-and-easy
tool for spelunking data until proven otherwise and/or working with small
datasets where the extra cost doesn't matter.

## See also

* https://sqlite.org/
* https://www.sqlite.org/json1.html
* https://www.gaia-gis.it/fossil/libspatialite/index
* https://github.com/whosonfirst/go-whosonfirst-sqlite
