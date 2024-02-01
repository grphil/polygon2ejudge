package import_contest

import (
	"fmt"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/lib/polygon_api"
	"strings"
)

func (t *ImportTask) ImportContest() {
	tmpDir, err := t.Transaction.GetTmp()
	if err != nil {
		return
	}

	contest, err := polygon_api.ImportContest(*t.PolygonContestId, tmpDir)
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("can not download polygon contest description, error: %s", err.Error()))
		return
	}

	fmt.Println("Imported contest description")

	err = t.Transaction.Commit("Imported contest description")
	if err != nil {
		return
	}

	for _, prob := range contest.Problems.Problems {
		shortName := strings.ToUpper(prob.Index)
		defaultEjudgeID := -1

		task := &import_problem.ImportTask{
			TaskCommon:     t.TaskCommon,
			PolygonProbUrl: &prob.Url,
			ShortName:      &shortName,
			EjudgeId:       &defaultEjudgeID,
			Transaction:    t.Transaction,
		}
		task.ImportProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}
}
