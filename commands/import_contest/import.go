package import_contest

import (
	"fmt"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/lib/polygon_api"
)

const englishLetters = 26

// It's a temporary solution, waiting for Polygon API changes.
func getShortName(id int) string {
	if id < englishLetters {
		return string('A' + id)
	}

	return "Z" + getShortName(id-englishLetters)
}

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

	for i, prob := range problems {
		shortName := getShortName(i)
		defaultEjudgeID := -1

		task := &import_problem.ImportTask{
			TaskCommon:    t.TaskCommon,
			PolygonProbID: &prob.ID,
			ShortName:     &shortName,
			EjudgeId:      &defaultEjudgeID,
			Transaction:   t.Transaction,
		}
		task.ImportProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}
}
