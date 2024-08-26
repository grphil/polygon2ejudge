package import_problem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (t *ImportTask) generateStatements() error {
	t.statement = &EjudgeStatementBuilder{
		Examples: make(map[int]*EjudgeExample),
	}

	err := t.extractPdfStatements()
	if err != nil {
		return fmt.Errorf("can not extract pdf statements, error: %s", err.Error())
	}

	if !*t.NoHtmlStatements {
		t.generateHTMLStatements()
	}

	xml, err := t.statement.GenerateXML()
	if err != nil {
		return fmt.Errorf("can not generate statements xml, error: %s", err.Error())
	}

	err = os.WriteFile(filepath.Join(t.ProbDir, "statements.xml"), xml, 0664)
	if err != nil {
		return fmt.Errorf("can not write statements xml, error: %s", err.Error())
	}

	t.config.Set("xml_file", "statements.xml")
	fmt.Println("Created statements")
	return nil
}

func (t *ImportTask) extractPdfStatements() error {
	var statements []*XProblemStatement

	if *t.AllLanguages {
		for _, statement := range t.problemXML.Statements.Statements {
			if statement.Type == "application/pdf" {
				statements = append(statements, statement)
			}
		}
	}

	if len(statements) == 0 {
		best := t.bestStatement("application/pdf")
		if best != nil {
			statements = append(statements, best)
		}
	}

	for _, statement := range statements {
		fileName := fmt.Sprintf("statements-%s.pdf", statement.Language)
		err := t.moveFile(statement.Path, "attachments", fileName)
		if err != nil {
			return err
		}
		statement.Path = fileName
		t.statement.PDFs = append(t.statement.PDFs, statement)
	}
	return nil
}

func (t *ImportTask) generateHTMLStatements() {
	err := t.extractAllFiles("statement-sections", "statement-sections")
	if err != nil {
		fmt.Printf("Warning: can not extract statement sections, error: %s\n", err.Error())
		return
	}

	s := t.bestStatement("application/x-tex")
	if s == nil {
		return
	}

	t.statement.Statement = &EjudgeStatement{}
	t.statement.Statement.Language = s.Language
	t.statementPath = filepath.Join(t.ProbDir, "statement-sections", s.Language)

	err = t.ProcessHTMLStatementsFiles()
	if err != nil {
		fmt.Printf("Warning: can not process html statement files, error: %s\n", err.Error())
	}
}

func (t *ImportTask) ProcessHTMLStatementsFiles() error {
	list, err := os.ReadDir(t.statementPath)
	if err != nil {
		return err
	}

	for _, file := range list {
		if file.IsDir() {
			continue
		}
		f := file.Name()
		if strings.HasSuffix(f, ".tex") {
			switch f {
			case "legend.tex":
				t.statement.Statement.Description = t.processTex(f)
			case "input.tex":
				t.statement.Statement.InputFormat = t.processTex(f)
			case "output.tex":
				t.statement.Statement.OutputFormat = t.processTex(f)
			case "interaction.tex":
				t.statement.Statement.Interaction = t.processTex(f)
			case "notes.tex":
				t.statement.Statement.Notes = t.processTex(f)
			case "scoring.tex":
				t.statement.Statement.Scoring = t.processTex(f)
			}
			continue
		}

		var exampleID int

		n, err := fmt.Sscanf(f, "example.%02d.a", &exampleID)
		if n == 1 && err == nil {
			t.addExample(exampleID, f, false)
			continue
		}

		n, err = fmt.Sscanf(f, "example.%02d", &exampleID)
		if n == 1 && err == nil {
			t.addExample(exampleID, f, true)
		}
	}

	return nil
}

func (t *ImportTask) addExample(id int, f string, input bool) {
	e, ok := t.statement.Examples[id]
	if !ok {
		e = &EjudgeExample{ID: id}
		t.statement.Examples[id] = e
	}

	text, err := os.ReadFile(filepath.Join(t.statementPath, f))
	if err != nil {
		fmt.Printf("Warning: can not read example %d file %s, error: %s\n", id, f, err.Error())
		return
	}
	if input {
		e.Input = string(text)
	} else {
		e.Output = string(text)
	}
}

func (t *ImportTask) bestStatement(statementType string) *XProblemStatement {
	var best *XProblemStatement
	for _, statement := range t.problemXML.Statements.Statements {
		if statement.Type != statementType {
			continue
		}

		if best == nil {
			best = statement
		}

		if best.Language == "russian" {
			break
		}

		if statement.Language == "russian" || statement.Language == "english" {
			best = statement
		}
	}
	return best
}
