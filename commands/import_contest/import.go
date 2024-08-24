package import_contest

import (
	"fmt"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/lib/polygon_api"
	"sort"
)

func (t *ImportTask) ImportContest() {
	problems, err := polygon_api.GetProblemsInContest(*t.PolygonContestId)
	if err != nil {
		t.Transaction.SetError(fmt.Errorf("can not get polygon contest problems, error: %s", err.Error()))
		return
	}

	fmt.Println("Imported problems list")

	err = t.Transaction.Commit("Imported problems list")
	if err != nil {
		return
	}
	shorts := make([]string, 0, len(problems))
	for short := range problems {
		shorts = append(shorts, short)
	}
	sort.Strings(shorts)

	for _, short := range shorts {
		short := short
		prob := problems[short]
		defaultEjudgeID := -1

		task := &import_problem.ImportTask{
			TaskCommon:    t.TaskCommon,
			PolygonProbID: &prob.ID,
			ShortName:     &short,
			EjudgeId:      &defaultEjudgeID,
			Transaction:   t.Transaction,
		}
		task.ImportProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}
}
