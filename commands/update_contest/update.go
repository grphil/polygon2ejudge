package update_contest

import "polygon2ejudge/commands/update_problem"

func (t *UpdateTask) UpdateContest() {
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

	for _, probID := range problems {
		task := &update_problem.UpdateTask{
			TaskCommon:      t.TaskCommon,
			EjudgeProblemId: &probID,
			Transaction:     t.Transaction,
		}
		task.UpdateProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}

	t.Transaction.Commit("Updated all problems in contest")
}
