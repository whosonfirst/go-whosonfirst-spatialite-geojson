package query

import (
	"database/sql"
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson"
)

func QueryToFeatureCollection(db *sql.DB, q string, args ...interface{}) (*geojson.FeatureCollection, error) {

	rows, err := db.Query(q, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return RowsToFeatureCollection(rows)
}

func RowsToFeatureCollection(rows *sql.Rows) (*geojson.FeatureCollection, error) {

	features := make([]*geojson.Feature, 0)

	for rows.Next() {

		var str_id string
		var str_props string
		var str_geom string

		err := rows.Scan(&str_id, &str_props, &str_geom)

		if err != nil {
			return nil, err
		}

		var props interface{}
		var geom interface{}

		err = json.Unmarshal([]byte(str_props), &props)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(str_geom), &geom)

		if err != nil {
			return nil, err
		}

		feature := geojson.Feature{
			Type:       "Feature",
			Id:         str_id,
			Properties: props,
			Geometry:   geom,
		}

		features = append(features, &feature)

	}

	feature_collection := geojson.FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	return &feature_collection, nil
}
