package submit_problem

import (
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/config"
)

func (t *SubmitTask) submitSolution(path string) error {
	var lang string

	switch filepath.Ext(path) {
	case ".cpp":
		lang = config.LangCpp
	case ".py":
		lang = config.LangPy
	case ".java":
		lang = config.LangJava
	case ".pas", ".dpr", ".fpc":
		lang = config.LangPas
	default:
		fmt.Printf("Unknown solution language for %s, skipping\n", filepath.Base(path))
		return nil
	}

	fileB, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can not read solution %s, error: %s", path, err.Error())
	}
	data := string(fileB)

	commentStart := "//"
	if lang == config.LangPy {
		commentStart = "#"
	}

	if *t.NoLint {
		data = fmt.Sprintf("%s %s\n%s", commentStart, config.UserConfig.NolintString, data)
	}

	data = fmt.Sprintf("%s %s\n%s", commentStart, filepath.Base(path), data)

	session, err := t.Transaction.GetEjudgeSession()
	if err != nil {
		return err
	}

	for _, langId := range config.GlobalConfig.LangIds[lang] {
		err = session.SubmitProblem(data, *t.ProblemId, langId)
		if err != nil {
			return fmt.Errorf("can not submit problem %d solution %s, error: %s",
				*t.ProblemId,
				filepath.Base(path),
				err.Error(),
			)
		}
	}
	return nil
}
