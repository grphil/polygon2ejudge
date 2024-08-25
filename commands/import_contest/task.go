package import_contest

import (
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	transaction2 "polygon2ejudge/lib/transaction"
)

type ImportTask struct {
	common.TaskCommon
	PolygonContestId *int

	Transaction *transaction2.Transaction
}

func AddImportContestCommand(parser *argparse.Parser) {
	task := &ImportTask{}
	ic := parser.AddCommand("ic", "Import contest from polygon", nil)
	task.AddCommonOptions(ic, true, false, true)
	task.PolygonContestId = ic.Int("", "polygon_id", &argparse.Option{
		Help:       "Polygon contest id",
		Required:   true,
		Positional: true,
	})

	ic.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.ImportContest()
		transaction.Finish()
	}
}
