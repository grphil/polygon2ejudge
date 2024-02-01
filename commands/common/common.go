package common

import (
	"github.com/hellflame/argparse"
)

type TaskCommon struct {
	ContestId *int

	// Submit options
	OnlyMain *bool
	NoLint   *bool

	// Import options
	NoOffline *bool

	Abstract        *string
	NoGenericParent *bool

	NoStatements     *bool
	NoHtmlStatements *bool
	AllLanguages     *bool

	NoCompileMainSolution *bool
	CompileAllSolutions   *bool

	NoTextareaInput *bool
	NoCustomRun     *bool
	NoPackageSave   *bool

	AllowFullReport      *bool
	FullReportSamplesAcm *bool

	NoConvertDprToPas *bool
}

func (t *TaskCommon) AddCommonOptions(parser *argparse.Parser, hasImport bool, hasSubmit bool) {
	t.ContestId = parser.Int("", "contest_id", &argparse.Option{
		Help:       "ID of ejudge contest",
		Required:   true,
		Positional: true,
	})

	if hasSubmit {
		t.OnlyMain = parser.Flag("m", "only-main", &argparse.Option{
			Help: "Submit only main solution (if specified)",
		})
		t.NoLint = parser.Flag("l", "no-lint", &argparse.Option{
			Help: "Add nolint string to solutions",
		})
	}

	if hasImport {
		t.Abstract = parser.String("", "abstract", &argparse.Option{
			Help: "Abstract \"parent\" problem to import all configs from",
		})
		t.NoGenericParent = parser.Flag("", "no-generic-parent", &argparse.Option{
			Help: "Do not add default generic parent if Generic problem is present in contest",
		})

		t.NoOffline = parser.Flag("", "no-offline", &argparse.Option{
			Help: "No offline groups for task",
		})

		t.NoStatements = parser.Flag("", "no-statements", &argparse.Option{
			Help: "Do not generate any statements (pdf and html)",
		})
		t.NoHtmlStatements = parser.Flag("", "no-html-statements", &argparse.Option{
			Help: "Do not generate html statements (pdf statements will be generated)",
		})
		t.AllLanguages = parser.Flag("l", "all-languages", &argparse.Option{
			Help: "Add pdf statements for all languages",
		})

		t.NoCompileMainSolution = parser.Flag("m", "no-compile-main", &argparse.Option{
			Help: "Do not compile main solution for the problem",
		})
		t.CompileAllSolutions = parser.Flag("", "compile-all", &argparse.Option{
			Help: "Compile all solutions for the problem",
		})

		t.NoTextareaInput = parser.Flag("", "no-textarea-input", &argparse.Option{
			Help: "Do not add textarea input for the problem",
		})
		t.NoCustomRun = parser.Flag("", "no-custom-run", &argparse.Option{
			Help: "Do not add custom code run support",
		})
		t.NoPackageSave = parser.Flag("", "no-package-save", &argparse.Option{
			Help: "Do not save package.zip for the problem",
		})

		t.AllowFullReport = parser.Flag("", "allow-full-report", &argparse.Option{
			Help: "Allow full report if group has complete report and is not sample",
		})
		t.FullReportSamplesAcm = parser.Flag("s", "full-report-samples-acm", &argparse.Option{
			Help: "Add full report for samples in acm problems",
		})

		t.NoConvertDprToPas = parser.Flag("", "no-dpr-to-pas", &argparse.Option{
			Help: "Do not change extension of all .dpr files to .pas",
		})
	}

}
