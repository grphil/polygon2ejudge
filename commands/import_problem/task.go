package import_problem

import (
	"archive/zip"
	"polygon2ejudge/commands/common"
	"polygon2ejudge/lib/orderedmap"
	"polygon2ejudge/lib/serve_cfg"
	transaction2 "polygon2ejudge/lib/transaction"

	"github.com/hellflame/argparse"
)

type ImportTask struct {
	common.TaskCommon
	PolygonProbUrl *string // unused now cause 403
	ShortName      *string
	EjudgeId       *int

	StatementsOnly bool

	PolygonProbID *int

	Transaction *transaction2.Transaction

	tmpDir      string
	packagePath string
	ProbDir     string

	InternalName string

	ServeCFG   *serve_cfg.ServeCFG
	problemXML *XProblemXML
	testset    *XTestset

	archive *zip.ReadCloser

	config            *orderedmap.OrderedMap
	problemOnlyConfig *orderedmap.OrderedMap

	deferFuncs []func()

	groups []*GroupInfo

	statement     *EjudgeStatementBuilder
	statementPath string
}

func AddImportProblemCommand(parser *argparse.Parser) {
	task := &ImportTask{}
	ip := parser.AddCommand("ip", "Import single problem from polygon", nil)
	task.AddCommonOptions(ip, true, false, true)
	task.PolygonProbID = ip.Int("", "problem_id", &argparse.Option{
		Help:       "Polygon problem ID",
		Required:   true,
		Positional: true,
	})
	task.ShortName = ip.String("", "short", &argparse.Option{Help: "Short name for the problem"})
	task.EjudgeId = ip.Int("", "ej-id", &argparse.Option{Help: "Ejudge id for the problem", Default: "-1"})

	ip.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.ImportProblem()
		transaction.Finish()
	}
}
