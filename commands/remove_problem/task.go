package remove_problem

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type RemoveTask struct {
	common.TaskCommon
	EjudgeProblemId *int

	Transaction *transaction2.Transaction
}

func AddRemoveProblemCommand(parser *argparse.Parser) {
	task := &RemoveTask{}
	rp := parser.AddCommand("rp", "Remove single problem from contest", nil)
	task.AddCommonOptions(rp, false)
	task.EjudgeProblemId = rp.Int("", "problem_id", &argparse.Option{
		Help:       "Ejudge id for the problem",
		Required:   true,
		Positional: true,
	})

	rp.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.RemoveProblem()
		transaction.Finish()
	}
}
