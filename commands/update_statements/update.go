package update_statements

import (
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/commands/import_problem"
)

func (t *UpdateTask) UpdateProblemStatements() {
	importTask, err := import_problem.CreateImportTaskForUpdate(t.TaskCommon, t.Transaction, *t.EjudgeProblemId)
	if err != nil {
		return
	}

	importTask.StatementsOnly = true
	importTask.ImportProblem()

	err = os.MkdirAll(filepath.Join(importTask.ProbDir, "attachments"), 0775)
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("error creating attachments directory, error: %v", err))
		return
	}

	err = t.Transaction.MovePath(
		filepath.Join(importTask.ProbDir, "statements.xml"),
		filepath.Join(importTask.ServeCFG.Path(), "problems", importTask.InternalName, "statements.xml"),
	)
	if err != nil {
		return
	}

	err = t.Transaction.MovePath(
		filepath.Join(importTask.ProbDir, "attachments"),
		filepath.Join(importTask.ServeCFG.Path(), "problems", importTask.InternalName, "attachments"),
	)
	if err != nil {
		return
	}

	t.Transaction.Commit(fmt.Sprintf("Updated statements form problem %s", *importTask.ShortName))
}
