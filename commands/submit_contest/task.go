package submit_contest

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type SubmitTask struct {
	common.TaskCommon

	Transaction *transaction2.Transaction
}

func AddSubmitContestCommand(parser *argparse.Parser) {
	task := &SubmitTask{}
	sc := parser.AddCommand("sc", "Submit solutions for all problems in contest", nil)
	task.AddCommonOptions(sc, false, true, false)

	sc.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		transaction.SetNoChange()
		task.Transaction = transaction
		task.SubmitContest()
		transaction.Finish()
	}
}
