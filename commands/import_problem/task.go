package import_problem

import (
	"archive/zip"
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/common"
	"polygon2ejudge/lib/orderedmap"
	"polygon2ejudge/lib/serve_cfg"
	transaction2 "polygon2ejudge/lib/transaction"
)

type ImportTask struct {
	common.TaskCommon
	PolygonProbId *int
	ShortName     *string
	EjudgeId      *int

	Transaction *transaction2.Transaction

	tmpDir      string
	packagePath string
	probDir     string

	internalName string

	serveCFG   *serve_cfg.ServeCFG
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
	ic := parser.AddCommand("ip", "Import single problem from polygon", nil)
	task.AddCommonOptions(ic, true)
	task.PolygonProbId = ic.Int("", "problem_id", &argparse.Option{
		Help:       "Polygon id for the problem",
		Required:   true,
		Positional: true,
	})
	task.ShortName = ic.String("", "short", &argparse.Option{Help: "Short name for the problem"})
	task.EjudgeId = ic.Int("", "ej-id", &argparse.Option{Help: "Ejudge id for the problem", Default: "-1"})

	ic.InvokeAction = func(invoked bool) {
		if !invoked {
			return
		}
		transaction := transaction2.NewTransaction(*task.ContestId)
		task.Transaction = transaction
		task.ImportProblem()
		transaction.Finish()
	}
}
