package index

import (
	"context"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	wof_index "github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	sql_index "github.com/whosonfirst/go-whosonfirst-sqlite/index"
	"io"
)

// THIS IS A TOTAL HACK UNTIL WE CAN SORT THINGS OUT IN
// go-whosonfirst-index... (20180206/thisisaaronland)

type Closer struct {
	fh io.Reader
}

func (c Closer) Read(b []byte) (int, error) {
	return c.fh.Read(b)
}

func (c Closer) Close() error {
	return nil
}

func NewDefaultSQLiteFeaturesIndexer(db sqlite.Database, to_index []sqlite.Table) (*sql_index.SQLiteIndexer, error) {

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) (interface{}, error) {

		select {

		case <-ctx.Done():
			return nil, nil
		default:
			path, err := wof_index.PathForContext(ctx)

			if err != nil {
				return nil, err
			}

			// HACK - see above
			closer := Closer{fh}

			i, err := feature.LoadFeatureFromReader(closer)

			if err != nil {
				msg := fmt.Sprintf("Unable to load %s, because %s", path, err)
				return nil, errors.New(msg)
			}

			return i, nil
		}
	}

	return sql_index.NewSQLiteIndexer(db, to_index, cb)
}
