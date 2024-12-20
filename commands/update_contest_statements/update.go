package update_contest_statements

import (
	"polygon2ejudge/commands/update_statements"
	"slices"
)

func (t *UpdateTask) UpdateContestStatements() {
	serveCfg, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	var problems []int

	for probID, prob := range serveCfg.Problems {
		if _, ok := prob.Get("extid"); ok {
			problems = append(problems, probID)
		}
	}

	slices.Sort(problems)

	for _, probID := range problems {
		task := &update_statements.UpdateTask{
			TaskCommon:      t.TaskCommon,
			EjudgeProblemId: &probID,
			Transaction:     t.Transaction,
		}
		task.UpdateProblemStatements()
		if t.Transaction.Err() != nil {
			return
		}
	}

	t.Transaction.Commit("Updated all problems in contest")
}
