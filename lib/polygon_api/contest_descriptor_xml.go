package polygon_api

type PContestXML struct {
	Names      PContestNamesXML      `xml:"names"`
	Statements PContestStatementsXML `xml:"statements"`
	Problems   PContestProblemsXML   `xml:"problems"`
}

type PContestNamesXML struct {
	Names []PContestNameXML `xml:"name"`
}

type PContestNameXML struct {
	Language string `xml:"language,attr"`
	Value    string `xml:"value,attr"`
}

type PContestStatementsXML struct {
	Statements []*PContestStatementXML `xml:"statement"`
}

type PContestStatementXML struct {
	Language string `xml:"language,attr"`
	Type     string `xml:"type,attr"`
	Url      string `xml:"url,attr"`
	Path     string `xml:"path,attr"`
}

type PContestProblemsXML struct {
	Problems []PContestProblemXML `xml:"problem"`
}

type PContestProblemXML struct {
	Index string `xml:"index,attr"`
	Url   string `xml:"url,attr"`
}
