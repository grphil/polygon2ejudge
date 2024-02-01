package update_contest

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type UpdateTask struct {
	common.TaskCommon

	Transaction *transaction2.Transaction
}

func AddUpdateContestCommand(parser *argparse.Parser) {
	task := &UpdateTask{}
	uc := parser.AddCommand("uc", "Update all tasks in ejudge contest", nil)
	task.AddCommonOptions(uc, true, false)

	uc.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.UpdateContest()
		transaction.Finish()
	}
}
