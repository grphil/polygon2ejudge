package update_contest_statements

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type UpdateTask struct {
	common.TaskCommon

	Transaction *transaction2.Transaction
}

func AddUpdateContestStatementsCommand(parser *argparse.Parser) {
	task := &UpdateTask{}
	ucs := parser.AddCommand("ucs", "Update all tasks statements in ejudge contest", nil)
	task.AddCommonOptions(ucs, false, false, true)

	ucs.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.UpdateContestStatements()
		transaction.Finish()
	}
}
