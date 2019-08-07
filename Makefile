# remember: if you're doing anything that _queries_ a database you'll need to pass 
# '--tags json1' to go build... (20180331/thisisaaronland)

tools:
	go build --mod vendor --tags json1 -o bin/wof-spatialite-geojson-index cmd/wof-spatialite-geojson-index/main.go
	go build --mod vendor --tags json1 -o bin/wof-spatialite-geojson-server cmd/wof-spatialite-geojson-server/main.go
