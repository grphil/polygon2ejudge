package update_problem

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type UpdateTask struct {
	common.TaskCommon

	EjudgeProblemId *int
	Transaction     *transaction2.Transaction
}

func AddUpdateProblemCommand(parser *argparse.Parser) {
	task := &UpdateTask{}
	up := parser.AddCommand("up", "Update single problem", nil)
	task.AddCommonOptions(up, true, false)
	task.EjudgeProblemId = up.Int("", "problem_id", &argparse.Option{
		Help:       "Ejudge id for the problem",
		Required:   true,
		Positional: true,
	})

	up.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.UpdateProblem()
		transaction.Finish()
	}
}
