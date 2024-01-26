package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"polygon2ejudge/commands/import_problem"
	"polygon2ejudge/commands/remove_problem"
	"polygon2ejudge/commands/reset"
)

func main() {
	parser := argparse.NewParser("polygon2ejudge", "Importer from polygon to ejudge", nil)

	import_problem.AddImportProblemCommand(parser)
	remove_problem.AddRemoveProblemCommand(parser)

	reset.AddResetConfigCommand(parser)

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}
}
