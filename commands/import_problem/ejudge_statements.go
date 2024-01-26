package import_problem

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"slices"
	"strings"
)

type EjudgeStatementBuilder struct {
	Statement *EjudgeStatement
	Examples  map[int]*EjudgeExample

	PDFs []*XProblemStatement
}

type EjudgeStatement struct {
	Language     string
	Description  string
	InputFormat  string
	OutputFormat string
	Interaction  string
	Notes        string
	Scoring      string
}

type EjudgeExample struct {
	Input  string `xml:"input"`
	Output string `xml:"output"`
	ID     int    `xml:"-"`
}

const kEjudgeStatementPrefix = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"no\"?>\n"

func (p *EjudgeStatementBuilder) GenerateXML() ([]byte, error) {
	xmlProblem := &XEjudgeProblemXML{}

	links := p.genPdfLinks()
	if p.Statement == nil && len(links) > 0 {
		p.Statement = &EjudgeStatement{}
	}

	if p.Statement != nil {
		p.Statement.Description = links + p.Statement.Description
		xmlProblem.Statements = p.Statement.genXmlStruct()
	}

	if len(p.Examples) > 0 {
		var examples []*EjudgeExample
		for _, e := range p.Examples {
			examples = append(examples, e)
		}
		slices.SortFunc(examples, func(a *EjudgeExample, b *EjudgeExample) int {
			return cmp.Compare(a.ID, b.ID)
		})
		xmlProblem.Examples = &XEjudgeExamples{examples}
	}

	resB, err := xml.Marshal(xmlProblem)
	if err != nil {
		return nil, err
	}
	res := kEjudgeStatementPrefix + string(resB)
	return []byte(res), nil
}

func (s *EjudgeStatement) genXmlStruct() *XEjudgeStatementXML {
	if len(s.Interaction) > 0 {
		s.Interaction = formatEjudgeTitle(kInteractionText, s.Language) + s.Interaction
		if len(s.InputFormat) > 0 {
			s.InputFormat = s.InputFormat + s.Interaction
		} else {
			s.Description = s.Description + s.Interaction
		}
	}
	if len(s.Scoring) > 0 {
		s.Scoring = formatEjudgeTitle(kScoringText, s.Language) + s.Scoring
		s.Notes = s.Notes + s.Scoring
	}

	xmlStatements := &XEjudgeStatementXML{}

	if len(s.Description) > 0 {
		xmlStatements.Description = &XMLVal{s.Description}
	}

	if len(s.InputFormat) > 0 {
		xmlStatements.InputFormat = &XMLVal{s.InputFormat}
	}

	if len(s.OutputFormat) > 0 {
		xmlStatements.OutputFormat = &XMLVal{s.OutputFormat}
	}

	if len(s.Notes) > 0 {
		xmlStatements.Notes = &XMLVal{s.Notes}
	}

	return xmlStatements
}

func (p *EjudgeStatementBuilder) genPdfLinks() string {
	links := strings.Builder{}

	for _, pdf := range p.PDFs {
		var linkText string
		if pdf.Language == "russian" {
			linkText = "Условие задачи (pdf)"
		} else {
			linkText = fmt.Sprintf("Statements %s (pdf)", pdf.Language)
		}
		link := fmt.Sprintf(
			"<p><a href=\"${getfile}=%s\">%s</a></p>",
			pdf.Path,
			linkText,
		)
		links.WriteString(link)
	}

	return links.String()
}

var kInteractionText = map[string]string{
	"russian": "<h3>Формат взаимодействия</h3>",
	"english": "<h3>Interaction</h3>",
}

var kScoringText = map[string]string{
	"russian": "<h3>Оценивание</h3>",
	"english": "<h3>Scoring</h3>",
}

func formatEjudgeTitle(m map[string]string, lang string) string {
	if lang == "russian" {
		return m["russian"]
	} else {
		return m["english"]
	}
}

type XEjudgeProblemXML struct {
	XMLName xml.Name `xml:"problem"`

	Statements *XEjudgeStatementXML `xml:"statement,omitempty"`
	Examples   *XEjudgeExamples     `xml:"examples,omitempty"`
}

type XEjudgeExamples struct {
	Examples []*EjudgeExample `xml:"example"`
}

type XEjudgeStatementXML struct {
	Description  *XMLVal `xml:"description,omitempty"`
	InputFormat  *XMLVal `xml:"input_format,omitempty"`
	OutputFormat *XMLVal `xml:"output_format,omitempty"`
	Notes        *XMLVal `xml:"notes,omitempty"`
}

type XMLVal struct {
	Val string `xml:",innerxml"`
}
