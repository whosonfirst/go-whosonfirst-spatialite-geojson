package http

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
	"strconv"
	"strings"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

type PropertiesFilter struct {
	Filters []string
	Args    []interface{}
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

func PropertiesFiltersFromQuery(req *gohttp.Request, param string) (*PropertiesFilter, error) {

	opts := sanitize.DefaultOptions()
	http_query := req.URL.Query()

	props, ok := http_query["property"]

	if !ok {
		return nil, nil
	}

	filters := make([]string, 0)
	args := make([]interface{}, 0)

	for _, raw_prop := range props {

		str_prop, err := sanitize.SanitizeString(raw_prop, opts)

		if err != nil {
			return nil, err
		}

		parts := strings.Split(str_prop, "=")

		if len(parts) != 2 {
			return nil, errors.New("Invalid property")
		}

		k := parts[0]
		v := parts[1]

		f := fmt.Sprintf("json_extract(properties, '$.%s') = ?", k)

		filters = append(filters, f)
		args = append(args, v)
	}

	p := PropertiesFilter{
		Filters: filters,
		Args:    args,
	}

	return &p, nil
}
