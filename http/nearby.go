package http

/*

curl -s 'localhost:9999/nearby?latitude=37.617342&longitude=-122.382932&property=wof:placetype%3Dvenue&radius=50' | jq '.features [].properties["wof:name"] '
"D-12 Wall Case"
"D59"
"Restroom Women's (Boarding Area D Terminal 2)"
"Every Beating Second "

*/

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-sqlite-geojson/query"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"log"
	gohttp "net/http"
	"strconv"
	"strings"
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

		radius := 200.0

		raw_radius := http_query.Get("radius")

		str_radius, err := sanitize.SanitizeString(raw_radius, opts)

		if err != nil {
			gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
			return
		}

		if str_radius != "" {

			fl_radius, err := strconv.ParseFloat(str_radius, 10)

			if err != nil {
				gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
				return
			}

			if fl_radius > 1000.0 || fl_radius < 0.0 {
				gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
				return
			}

			radius = fl_radius
		}

		// because this
		// https://stackoverflow.com/questions/8287769/what-unit-is-does-spatialites-distance-function-return

		radius = radius / 111120.0

		pt := fmt.Sprintf("POINT(%0.6f %0.6f)", coord.Longitude, coord.Latitude)

		q := fmt.Sprintf(`SELECT id, properties, AsGeoJSON(geometry) FROM geojson WHERE
			PtDistWithin(
				ST_GeomFromText('%s'),
				ST_Centroid(geometry),
				%f)`, pt, radius)

		args := make([]interface{}, 0)

		props, ok := http_query["property"]

		if ok {

			for _, raw_prop := range props {

				str_prop, err := sanitize.SanitizeString(raw_prop, opts)

				if err != nil {
					gohttp.Error(rsp, "Invalid property", gohttp.StatusBadRequest)
					return
				}

				parts := strings.Split(str_prop, "=")

				if len(parts) != 2 {
					gohttp.Error(rsp, "Invalid property", gohttp.StatusBadRequest)
					return
				}

				k := parts[0]
				v := parts[1]

				q = fmt.Sprintf("%s AND json_extract(properties, '$.%s') = ?", q, k)
				args = append(args, v)
			}
		}

		conn, err := db.Conn()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		log.Println(q, args)

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
