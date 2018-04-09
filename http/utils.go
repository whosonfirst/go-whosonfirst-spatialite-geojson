package http

import (
	"errors"
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
	"strconv"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

func CoordFromQuery(req *gohttp.Request) (*Coord, error) {

	lat, err := LatitudeFromQuery(req, "latitude")

	if err != nil {
		return nil, err
	}

	lon, err := LongitudeFromQuery(req, "longitude")

	if err != nil {
		return nil, err
	}

	c := Coord{
		Latitude:  lat,
		Longitude: lon,
	}

	return &c, nil
}

func LatitudeFromQuery(req *gohttp.Request, param string) (float64, error) {

	opts := sanitize.DefaultOptions()
	http_query := req.URL.Query()

	raw_lat := http_query.Get(param)

	str_lat, err := sanitize.SanitizeString(raw_lat, opts)

	if err != nil {
		return 0.0, errors.New("Invalid latitude string")
	}

	if str_lat == "" {
		return 0.0, errors.New("Empty latitude string")
	}

	lat, err := strconv.ParseFloat(str_lat, 10)

	if err != nil {
		return 0.0, err
	}

	if lat > 90.0 || lat < -90.0 {
		return 0.0, err
	}

	return lat, nil
}

func LongitudeFromQuery(req *gohttp.Request, param string) (float64, error) {

	opts := sanitize.DefaultOptions()
	http_query := req.URL.Query()

	raw_lon := http_query.Get(param)

	str_lon, err := sanitize.SanitizeString(raw_lon, opts)

	if err != nil {
		return 0.0, errors.New("Invalid longitude string")
	}

	if str_lon == "" {
		return 0.0, errors.New("Empty longitude string")
	}

	lon, err := strconv.ParseFloat(str_lon, 10)

	if err != nil {
		return 0.0, err
	}

	if lon > 180.0 || lon < -180.0 {
		return 0.0, err
	}

	return lon, nil
}
