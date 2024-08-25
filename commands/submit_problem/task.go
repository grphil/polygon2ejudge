package submit_problem

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type SubmitTask struct {
	common.TaskCommon
	ProblemId *int

	Transaction *transaction2.Transaction
}

func AddSubmitProblemCommand(parser *argparse.Parser) {
	task := &SubmitTask{}
	sp := parser.AddCommand("sp", "Submit solutions for single problem", nil)
	task.AddCommonOptions(sp, false, true, false)

	task.ProblemId = sp.Int("", "problem_id", &argparse.Option{
		Help:       "Ejudge id for the problem",
		Required:   true,
		Positional: true,
	})

	sp.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		transaction.SetNoChange()
		task.Transaction = transaction
		task.SubmitProblem()
		transaction.Finish()
	}
}
