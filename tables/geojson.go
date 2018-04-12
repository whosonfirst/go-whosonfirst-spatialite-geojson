package tables

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/twpayne/go-geom"
	gogeom_geojson "github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	wof_geom "github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-spatialite-geojson"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	_ "log"
)

type GeoJSONTable struct {
	geojson.GeoJSONTable
	name string
}

func NewGeoJSONTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewGeoJSONTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewGeoJSONTable() (sqlite.Table, error) {

	t := GeoJSONTable{
		name: "geojson",
	}

	return &t, nil
}

func (t *GeoJSONTable) Name() string {
	return t.name
}

func (t *GeoJSONTable) Schema() string {

	sql := `CREATE TABLE %s (
		id TEXT NOT NULL PRIMARY KEY,
		properties JSON
	);

	SELECT InitSpatialMetaData();
	SELECT AddGeometryColumn('%s', 'geometry', 4326, 'GEOMETRY', 'XY');
	SELECT CreateSpatialIndex('%s', 'geometry');
	`

	return fmt.Sprintf(sql, t.Name(), t.Name(), t.Name())
}

func (t *GeoJSONTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *GeoJSONTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexFeature(db, i.(wof_geojson.Feature))
}

func (t *GeoJSONTable) IndexFeature(db sqlite.Database, f wof_geojson.Feature) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	str_id := f.Id()

	str_geom, err := wof_geom.ToString(f)

	if err != nil {
		return err
	}

	// but wait! there's more!! for reasons I've forgotten (simonw told me)
	// the spatialite doesn't really like indexing GeomFromGeoJSON but also
	// doesn't complain about it - it just chugs along happily filling your
	// database with null geometries so we're going to take advantage of the
	// handy "go-geom" package to convert the GeoJSON geometry in to WKT -
	// it is "one more thing" to import and maybe it would be better to just
	// write a custom converter but not today...
	// (20180122/thisisaaronland)

	var g geom.T
	err = gogeom_geojson.Unmarshal([]byte(str_geom), &g)

	if err != nil {
		return err
	}

	str_wkt, err := wkt.Marshal(g)

	if err != nil {
		return err
	}

	props := gjson.GetBytes(f.Bytes(), "properties")

	if !props.Exists() {
		return errors.New("feature is missing properties")
	}

	b_props, err := json.Marshal(props.Value())

	if err != nil {
		return err
	}

	str_props := string(b_props)

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		id, properties, geometry
	) VALUES (
		?, ?, GeomFromText('%s', 4326)
	)`, t.Name(), str_wkt)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(str_id, str_props)

	if err != nil {
		return err
	}

	return tx.Commit()
}
