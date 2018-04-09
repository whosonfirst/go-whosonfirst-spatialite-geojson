package http

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

		opts := sanitize.DefaultOptions()

		http_query := req.URL.Query()

		raw_lat := http_query.Get("latitude")
		raw_lon := http_query.Get("longitude")

		str_lat, err := sanitize.SanitizeString(raw_lat, opts)

		if err != nil {
			gohttp.Error(rsp, "Invalid latitude", gohttp.StatusBadRequest)
			return
		}

		if str_lat == "" {
			gohttp.Error(rsp, "Missing latitude", gohttp.StatusBadRequest)
			return
		}

		str_lon, err := sanitize.SanitizeString(raw_lon, opts)

		if err != nil {
			gohttp.Error(rsp, "Invalid longitude", gohttp.StatusBadRequest)
			return
		}

		if str_lon == "" {
			gohttp.Error(rsp, "Missing longitude", gohttp.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(str_lat, 10)

		if err != nil {
			gohttp.Error(rsp, "Invalid latitude", gohttp.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(str_lon, 10)

		if err != nil {
			gohttp.Error(rsp, "Invalid longitude", gohttp.StatusBadRequest)
			return
		}

		radius := 200

		raw_radius := http_query.Get("radius")

		str_radius, err := sanitize.SanitizeString(raw_radius, opts)

		if err != nil {
			gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
			return
		}

		if str_radius != "" {

			int_radius, err := strconv.Atoi(str_radius)

			if err != nil {
				gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
				return
			}

			if int_radius > 1000 || int_radius < 0 {
				gohttp.Error(rsp, "Invalid radius", gohttp.StatusBadRequest)
				return
			}

			radius = int_radius
		}

		pt := fmt.Sprintf("POINT(%0.6f %0.6f)", lon, lat)

		q := fmt.Sprintf(`SELECT id, properties, AsGeoJSON(geometry) FROM geojson WHERE
			PtDistWithin(
				ST_Centroid(geometry),
				ST_GeomFromText('%s'),
				%d)`, pt, radius)

		args := make([]interface{}, 0)

		props, ok := http_query["property"]

		if ok {

			log.Println("PROPS", props)

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
				log.Println("Q", q)

				args = append(args, v)
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
