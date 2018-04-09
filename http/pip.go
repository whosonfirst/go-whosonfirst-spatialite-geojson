package http

import (
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-sqlite-geojson/query"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"log"
	gohttp "net/http"
)

func PointInPolygonHandler(db *database.SQLiteDatabase) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		coord, err := CoordFromQuery(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		pt := fmt.Sprintf("POINT(%0.6f %0.6f)", coord.Longitude, coord.Latitude)

		q := fmt.Sprintf(`SELECT id, properties, AsGeoJSON(geometry) FROM geojson WHERE
			Contains(
				geometry,
				ST_GeomFromText('%s')
				)`, pt)

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
