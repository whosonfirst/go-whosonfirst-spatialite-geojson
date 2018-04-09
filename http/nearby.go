package http

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-sqlite-geojson/query"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"log"
	gohttp "net/http"
	"strconv"
)

func NearbyHandler(db *database.SQLiteDatabase) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		http_query := req.URL.Query()

		str_lat := http_query.Get("latitude")
		str_lon := http_query.Get("longitude")

		if str_lat == "" {
			gohttp.Error(rsp, "Missing latitude", gohttp.StatusBadRequest)
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

		conn, err := db.Conn()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		pt := fmt.Sprintf("POINT(%0.6f %0.6f)", lon, lat)

		q := fmt.Sprintf(`SELECT id, properties, AsGeoJSON(geometry) FROM geojson WHERE
			PtDistWithin(
				ST_GeomFromText('%s'),
				ST_Centroid(geometry),
				200)`, pt)

		log.Println(q)
				
		fc, err := query.QueryToFeatureCollection(conn, q)

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
