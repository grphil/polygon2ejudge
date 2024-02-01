package import_problem

type XProblemXML struct {
	ShortName string `xml:"short-name,attr"`
	Revision  string `xml:"revision,attr"`
	Url       string `xml:"url,attr"`

	Names      XProblemNames      `xml:"names"`
	Statements XProblemStatements `xml:"statements"`

	Judging XJudging `xml:"judging"`
	Files   XFiles   `xml:"files"`
	Assets  XAssets  `xml:"assets"`
}

type XProblemNames struct {
	Names []*XProblemName `xml:"name"`
}

type XProblemName struct {
	Language string `xml:"language,attr"`
	Value    string `xml:"value,attr"`
}

type XProblemStatements struct {
	Statements []*XProblemStatement `xml:"statement"`
}

type XProblemStatement struct {
	Language string `xml:"language,attr"`
	Type     string `xml:"type,attr"`
	Path     string `xml:"path,attr"`
}

type XJudging struct {
	InputFile  string      `xml:"input-file,attr"`
	OutputFile string      `xml:"output-file,attr"`
	Testsets   []*XTestset `xml:"testset"`
}

type XTestset struct {
	Name        string  `xml:"name,attr"`
	TimeLimit   int     `xml:"time-limit"`
	MemoryLimit int     `xml:"memory-limit"`
	Tests       XTests  `xml:"tests"`
	Groups      XGroups `xml:"groups"`
}

type XTests struct {
	Tests []*XTest `xml:"test"`
}

type XTest struct {
	Group  string  `xml:"group,attr,omitempty"`
	Points float64 `xml:"points,attr,omitempty"`
	Sample bool    `xml:"sample,attr,omitempty"`
}

type XGroups struct {
	Groups []*XGroup `xml:"group"`
}

type XGroup struct {
	Name           string             `xml:"name,attr"`
	FeedbackPolicy string             `xml:"feedback-policy,attr"`
	PointsPolicy   string             `xml:"points-policy,attr"`
	Points         float64            `xml:"points,attr"`
	Dependencies   XGroupDependencies `xml:"dependencies"`
}

type XGroupDependencies struct {
	Dependencies []*XGroupDependency `xml:"dependency"`
}

type XGroupDependency struct {
	Group string `xml:"group,attr"`
}

type XFiles struct {
	Resources XResources `xml:"resources"`
}

type XResources struct {
	Resources []*XSource `xml:"file"`
}

type XAssets struct {
	Checker    XSourceFile  `xml:"checker"`
	Interactor *XSourceFile `xml:"interactor,omitempty"`
	Solutions  XSolutions   `xml:"solutions"`
}

type XSourceFile struct {
	Tag    string  `xml:"tag,attr,omitempty"`
	Source XSource `xml:"source"`
}

type XSource struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type XSolutions struct {
	Solutions []*XSourceFile `xml:"solution"`
}
