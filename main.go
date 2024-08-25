package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/import_contest"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/commands/remove_contest"
	"polygon2ejudge/commands/remove_problem"
	"polygon2ejudge/commands/reset"
	"polygon2ejudge/commands/submit_contest"
	"polygon2ejudge/commands/submit_problem"
	"polygon2ejudge/commands/update_contest"
	"polygon2ejudge/commands/update_problem"
	"polygon2ejudge/commands/update_statements"
)

func main() {
	parser := argparse.NewParser("polygon2ejudge", "Importer from polygon to ejudge", nil)

	import_problem.AddImportProblemCommand(parser)
	import_contest.AddImportContestCommand(parser)

	update_problem.AddUpdateProblemCommand(parser)
	update_contest.AddUpdateContestCommand(parser)

	update_statements.AddUpdateStatementsCommand(parser)

	remove_problem.AddRemoveProblemCommand(parser)
	remove_contest.AddRemoveContestCommand(parser)

	submit_problem.AddSubmitProblemCommand(parser)
	submit_contest.AddSubmitContestCommand(parser)

	reset.AddResetConfigCommand(parser)

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
}
