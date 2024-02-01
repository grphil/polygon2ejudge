package submit_problem

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (t *SubmitTask) SubmitProblem() {
	serveCFG, err := t.Transaction.EditServeCfg()
	if err != nil {
		return
	}

	prob, ok := serveCFG.Problems[*t.ProblemId]
	if !ok {
		t.Transaction.SetError(fmt.Errorf("problem with id %d not found", *t.ProblemId))
		return
	}

	internalName := prob.GetStr("internal_name")
	if len(internalName) == 0 {
		t.Transaction.SetError(fmt.Errorf("no internal name for problem %d", *t.ProblemId))
		return
	}

	probDir := filepath.Join(serveCFG.Path(), "problems", internalName)

	mainSolution := prob.GetStr("solution_cmd")
	if len(mainSolution) > 0 {
		files, err := os.ReadDir(probDir)
		if err != nil {
			t.Transaction.SetError(fmt.Errorf("can not read problem files, error: %s", err.Error()))
			return
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			if strings.HasPrefix(f.Name(), mainSolution) {
				err = t.submitSolution(filepath.Join(probDir, f.Name()))
				if err != nil {
					t.Transaction.SetError(err)
					return
				}
			}
		}
	}

	if !*t.OnlyMain {
		solutionsDir, ok := getSolutionsDir(probDir)
		if ok {
			files, err := os.ReadDir(solutionsDir)
			if err != nil {
				t.Transaction.SetError(fmt.Errorf("can not list solutions for problem %d, error: %s", *t.ProblemId, err.Error()))
				return
			}

			slices.SortFunc(files, func(a os.DirEntry, b os.DirEntry) int {
				return cmp.Compare(a.Name(), b.Name())
			})

			for _, f := range files {
				if f.IsDir() {
					continue
				}

				err = t.submitSolution(filepath.Join(solutionsDir, f.Name()))
				if err != nil {
					t.Transaction.SetError(err)
					return
				}
			}
		}
	}

	t.Transaction.Commit(fmt.Sprintf("Submitted solutions for problem %d %s", *t.ProblemId, internalName))
}

func getSolutionsDir(probDir string) (string, bool) {
	solutionsDir := filepath.Join(probDir, "solutions")
	if fileExists(solutionsDir) {
		return solutionsDir, true
	}
	solutionsDir = filepath.Join(probDir, "solutions1")
	return solutionsDir, fileExists(solutionsDir)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
