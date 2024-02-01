package remove_contest

import (
	"polygon2ejudge/commands/remove_problem"
	"slices"
)

func (t *RemoveTask) RemoveContest() {
	serveCFG, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	var problems []int

	for probID, _ := range serveCFG.Problems {
		problems = append(problems, probID)
	}

	slices.Sort(problems)

	for _, probID := range problems {
		task := remove_problem.RemoveTask{
			TaskCommon:      t.TaskCommon,
			EjudgeProblemId: &probID,
			Transaction:     t.Transaction,
			KeepServeCfg:    false,
		}
		task.RemoveProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}
	t.Transaction.Commit("Removed all problems from contest")
}
