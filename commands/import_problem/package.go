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

	fmt.Printf("downloading package for problem %d\n", *t.PolygonProbId)
	err = polygon_api.ImportPackage(*t.PolygonProbId, t.packagePath)
	if err != nil {
		return fmt.Errorf("polygon api error while loading package: %s", err.Error())
	}

	return nil
}
