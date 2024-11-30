package import_problem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/polygon_api"
)

func (t *ImportTask) importPackage() error {
	t.packagePath = filepath.Join(t.tmpDir, "package", "package.zip")
	err := os.MkdirAll(filepath.Dir(t.packagePath), 0775)
	if err != nil {
		return err
	}

	fmt.Println("downloading problem package")
	if t.PolygonProbID != nil {
		err = polygon_api.ImportPackageApi(*t.PolygonProbID, t.packagePath)
	} else if t.PolygonProbUrl != nil { // TODO: wait for fix from Mike
		err = polygon_api.ImportPackage(*t.PolygonProbUrl, t.packagePath)
	} else {
		return errors.New("no problem id or url provided")
	}

	if err != nil {
		return fmt.Errorf("polygon error while loading package: %s", err.Error())
	}

	return nil
}
