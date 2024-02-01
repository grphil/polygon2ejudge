package remove_contest

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type RemoveTask struct {
	common.TaskCommon

	Transaction *transaction2.Transaction
}

func AddRemoveContestCommand(parser *argparse.Parser) {
	task := &RemoveTask{}
	rc := parser.AddCommand("rc", "Remove all problems from ejudge contest", nil)
	task.AddCommonOptions(rc, false, false)

	rc.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.RemoveContest()
		transaction.Finish()
	}
}
