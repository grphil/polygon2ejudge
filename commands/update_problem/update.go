package update_problem

import (
	"fmt"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/commands/remove_problem"
	"strings"
)

func (t *UpdateTask) UpdateProblem() {
	serveCfg, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	prob, ok := serveCfg.Problems[*t.EjudgeProblemId]
	if !ok {
		t.Transaction.SetError(fmt.Errorf("problem %d not found", *t.EjudgeProblemId))
		return
	}

	extID := prob.GetStr("extid")

	var polygonID int
	var problemUrl string

	if strings.HasPrefix(extID, "https://polygon.codeforces.com/") {
		problemUrl = extID
	} else {
		n, err := fmt.Sscanf(extID, "polygon:%d", &polygonID)
		if err != nil {
			t.Transaction.SetError(fmt.Errorf("can not get problem %d polygon id, error: %s", *t.EjudgeProblemId, err.Error()))
			return
		}
		if n == 0 {
			t.Transaction.SetError(fmt.Errorf("can not get problem %d polygon id", *t.EjudgeProblemId))
			return
		}
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

	shortName := prob.GetStr("short_name")
	importTask := &import_problem.ImportTask{
		TaskCommon:  t.TaskCommon,
		ShortName:   &shortName,
		EjudgeId:    t.EjudgeProblemId,
		Transaction: t.Transaction,
	}

	if len(problemUrl) > 0 {
		importTask.PolygonProbUrl = &problemUrl
	} else {
		importTask.PolygonProbID = &polygonID
	}
	importTask.ImportProblem()
}
