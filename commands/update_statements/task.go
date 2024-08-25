package update_statements

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

func AddUpdateStatementsCommand(parser *argparse.Parser) {
	task := &UpdateTask{}
	us := parser.AddCommand("us", "Update problem statements", nil)
	task.AddCommonOptions(us, false, false, true)
	task.EjudgeProblemId = us.Int("", "problem_id", &argparse.Option{
		Help:       "Ejudge id for the problem",
		Required:   true,
		Positional: true,
	})

	us.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.UpdateProblemStatements()
		transaction.Finish()
	}
}
