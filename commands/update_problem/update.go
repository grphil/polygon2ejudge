package update_problem

import (
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/commands/remove_problem"
)

func (t *UpdateTask) UpdateProblem() {
	importTask, err := import_problem.CreateImportTaskForUpdate(t.TaskCommon, t.Transaction, *t.EjudgeProblemId)
	if err != nil {
		return
	}

	removeTask := &remove_problem.RemoveTask{
		TaskCommon:      t.TaskCommon,
		EjudgeProblemId: t.EjudgeProblemId,
		Transaction:     t.Transaction,
		KeepServeCfg:    true,
	}
	removeTask.RemoveProblem()
	if t.Transaction.Err() != nil {
		return
	}

	importTask.ImportProblem()
}
