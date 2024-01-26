package remove_problem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (t *RemoveTask) RemoveProblem() {
	serveCFG, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	prob, ok := serveCFG.Problems[*t.EjudgeProblemId]
	if !ok {
		t.Transaction.SetError(fmt.Errorf("can not remove problem with id %d, no problem found", *t.EjudgeProblemId))
		return
	}

	internalName := prob.GetStr("internal_name")
	if len(internalName) == 0 {
		t.Transaction.SetError(fmt.Errorf("can not remove problem with id %d, problem has no dir", *t.EjudgeProblemId))
		return
	}

	packageZipPath := filepath.Join(serveCFG.Path(), "download", fmt.Sprintf("%s.zip", internalName))

	_, err = os.Stat(packageZipPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Transaction.SetError(err)
			return
		}
	} else {
		err = t.Transaction.RemovePath(packageZipPath)
		if err != nil {
			return
		}
	}

	delete(serveCFG.Problems, *t.EjudgeProblemId)
	err = t.Transaction.RemovePath(filepath.Join(serveCFG.Path(), "problems", internalName))
	if err != nil {
		return
	}
	t.Transaction.Commit(fmt.Sprintf("Removed problem id: %d name: %s", *t.EjudgeProblemId, internalName))
}
