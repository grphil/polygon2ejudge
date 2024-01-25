package import_problem

import (
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/polygon_api"
)

func (t *ImportTask) importPackage() error {
	t.packagePath = filepath.Join(t.tmpDir, "package", "package.zip")
	err := os.MkdirAll(filepath.Dir(t.packagePath), 0774)
	if err != nil {
		return err
	}

	if t.PolygonApi == nil {
		t.PolygonApi, err = polygon_api.NewPolygonApi(*t.ResetCredentials)
		if err != nil {
			return fmt.Errorf("can not create polygon api, err: %s", err.Error())
		}
	}

	err = t.PolygonApi.ImportPackage(*t.PolygonProbId, t.packagePath)
	if err != nil {
		return fmt.Errorf("polygon api error while loading package: %s", err.Error())
	}

	err = t.Transaction.MovePath(t.packagePath, filepath.Join(t.serveCFG.Path(), "download", "pack.zip"))
	if err != nil {
		return err
	}

	return nil
}
