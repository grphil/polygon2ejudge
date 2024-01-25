package polygon_api

import (
	"polygon2ejudge/lib/config"
)

const POLYGON_URL = "https://polygon.codeforces.com/api/"

type PolygonApi struct {
	confFile *config.ConfFile
}

func NewPolygonApi(resetCredentials bool) (*PolygonApi, error) {
	confFile, err := config.GetConfigFile(resetCredentials)
	if err != nil {
		return nil, err
	}
	return &PolygonApi{confFile}, err
}
