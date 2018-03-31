package geojson

import (
       "github.com/whosonfirst/go-whosonfirst-sqlite"
       "github.com/whosonfirst/go-whosonfirst-geojson-v2"
)

type GeoJSONTable interface {
     sqlite.Table
     IndexGeoJSON(sqlite.Database, geojson.Feature) error
}
