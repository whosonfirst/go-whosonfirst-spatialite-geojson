package geojson

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
)

type GeoJSONTable interface {
	sqlite.Table
	IndexGeoJSON(sqlite.Database, geojson.Feature) error
}

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

type Feature struct {
	Type       string      `json:"type"`
	Id         string      `json:"id"`
	Properties interface{} `json:"properties"`
	Geometry   interface{} `json:"geometry"`
}
