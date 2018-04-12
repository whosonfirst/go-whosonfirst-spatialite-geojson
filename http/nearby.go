package http

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson/query"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	_ "log"
	gohttp "net/http"
	"strconv"
)

func NearbyHandler(db *database.SQLiteDatabase) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		coord, err := CoordFromQuery(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		opts := sanitize.DefaultOptions()

		http_query := req.URL.Query()

		distance := 50.0

		raw_distance := http_query.Get("distance")

		str_distance, err := sanitize.SanitizeString(raw_distance, opts)

		if err != nil {
			gohttp.Error(rsp, "Invalid distance", gohttp.StatusBadRequest)
			return
		}

		if str_distance != "" {

			fl_distance, err := strconv.ParseFloat(str_distance, 10)

			if err != nil {
				gohttp.Error(rsp, "Invalid distance", gohttp.StatusBadRequest)
				return
			}

			if fl_distance > 1000.0 || fl_distance < 0.0 {
				gohttp.Error(rsp, "Invalid distance", gohttp.StatusBadRequest)
				return
			}

			distance = fl_distance
		}

		// because this
		// https://stackoverflow.com/questions/8287769/what-unit-is-does-spatialites-distanceance-function-return

		distance = distance / 111120.0

		pt := fmt.Sprintf("POINT(%0.6f %0.6f)", coord.Longitude, coord.Latitude)

		q := fmt.Sprintf(`SELECT id, properties, AsGeoJSON(geometry) FROM geojson WHERE
			PtDistWithin(
				ST_GeomFromText('%s'),
				ST_Centroid(geometry),
				%f)`, pt, distance)

		args := make([]interface{}, 0)

		filters, err := PropertiesFiltersFromQuery(req, "properties")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		if filters != nil {

			for _, f := range filters.Filters {
				q = fmt.Sprintf("%s AND %s", q, f)
			}

			for _, a := range filters.Args {
				args = append(args, a)
			}
		}

		conn, err := db.Conn()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		fc, err := query.QueryToFeatureCollection(conn, q, args...)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(fc)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(js)
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
