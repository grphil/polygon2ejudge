package submit_contest

import (
	"polygon2ejudge/commands/submit_problem"
)

func (t *SubmitTask) SubmitContest() {
	serveCFG, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	var problems []int

	for probID, _ := range serveCFG.Problems {
		problems = append(problems, probID)
	}

	for _, probID := range problems {
		task := &submit_problem.SubmitTask{
			TaskCommon:  t.TaskCommon,
			ProblemId:   &probID,
			Transaction: t.Transaction,
		}
		task.SubmitProblem()
		if t.Transaction.Err() != nil {
			return
		}
	}
	t.Transaction.Commit("Submitted all problems in contest")
}
