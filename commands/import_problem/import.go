package import_problem

import (
	"fmt"
	"path/filepath"
)

func (t *ImportTask) ImportProblem() {
	defer func() {
		for i := len(t.deferFuncs) - 1; i >= 0; i-- {
			t.deferFuncs[i]()
		}
	}()

	t.addDeferFunc(func() {
		t.Transaction.Commit(fmt.Sprintf("Imported problem %d name: %s", *t.PolygonProbId, t.internalName))
	})

	var err error
	t.tmpDir, err = t.Transaction.GetTmp()
	if err != nil {
		return
	}

	t.serveCFG, err = t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	err = t.importPackage()
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("error while loading package, error: %s", err.Error()))
		return
	}

	err = t.extractProblemFiles()
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("error while extracting zip package, error: %s", err.Error()))
		return
	}

	err = t.fillInConfig()
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("error while building problem cfg, error: %s", err.Error()))
		return
	}

	err = t.buildValuer()
	if err != nil {
		t.Transaction.SetError(err)
		return
	}

	if !*t.NoStatements {
		err = t.generateStatements()
		if err != nil {
			t.Transaction.SetError(err)
			return
		}
	}

	t.Transaction.MovePath(
		t.probDir,
		filepath.Join(t.serveCFG.Path(), "problems", t.internalName),
	)
}

func (t *ImportTask) addDeferFunc(f func()) {
	t.deferFuncs = append(t.deferFuncs, f)
}
