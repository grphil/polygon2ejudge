package serve_cfg

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/config"
	"polygon2ejudge/lib/orderedmap"
	"strconv"
	"strings"
)

// We can not use default go toml libs, as they lack comments and unquoted strings, that are used in ejudge

type ServeCFG struct {
	Global    *orderedmap.OrderedMap
	languages []*orderedmap.OrderedMap
	testers   []*orderedmap.OrderedMap
	other     map[string][]*orderedmap.OrderedMap

	Problems         map[int]*orderedmap.OrderedMap
	abstractProblems []*orderedmap.OrderedMap

	ContestID int

	HasGeneric   bool
	MaxProblemId int
}

func NewServeCFG(contestID int) (*ServeCFG, error) {
	s := &ServeCFG{
		Global:    orderedmap.New(),
		languages: make([]*orderedmap.OrderedMap, 0),
		testers:   make([]*orderedmap.OrderedMap, 0),
		other:     make(map[string][]*orderedmap.OrderedMap),

		Problems:         make(map[int]*orderedmap.OrderedMap),
		abstractProblems: make([]*orderedmap.OrderedMap, 0),

		ContestID: contestID,

		HasGeneric:   false,
		MaxProblemId: 0,
	}

	cfgFile, err := os.Open(filepath.Join(s.Path(), "conf", "serve.cfg"))
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	currSection := orderedmap.New()
	currSectionName := "Global"

	var comments string

	scanner := bufio.NewScanner(cfgFile)
	for scanner.Scan() {
		initialTxt := scanner.Text()
		txt := strings.TrimSpace(initialTxt)

		if len(txt) == 0 {
			continue
		}

		if strings.HasPrefix(txt, "#") {
			comments = comments + "\n" + initialTxt
			continue
		}

		if len(comments) > 0 {
			currSection.Set(comments, orderedmap.Comment(true))
			comments = ""
		}

		if strings.HasPrefix(txt, "[") && strings.HasSuffix(txt, "]") {
			s.addSection(currSectionName, currSection)
			currSection = orderedmap.New()
			currSectionName = txt[1 : len(txt)-1]
			continue
		}

		if strings.Contains(txt, "=") {
			key := strings.TrimSpace(txt[:strings.Index(txt, "=")])
			val := strings.TrimSpace(txt[strings.Index(txt, "=")+1:])

			if val == "true" {
				currSection.Set(key, true)
			} else if val == "false" {
				currSection.Set(key, false)
			} else if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
				currSection.Set(key, val[1:len(val)-1])
			} else if valInt, err := strconv.Atoi(val); err == nil {
				currSection.Set(key, valInt)
			} else {
				currSection.Set(key, orderedmap.UnquotedStr(val))
			}
		} else {
			currSection.Set(txt, orderedmap.NoVal(true))
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	if len(comments) > 0 {
		currSection.Set(comments, orderedmap.Comment(true))
	}
	s.addSection(currSectionName, currSection)
	return s, nil
}

func (s *ServeCFG) Problem(id int) *orderedmap.OrderedMap {
	prob, ok := s.Problems[id]
	if ok {
		return prob
	}

	prob = orderedmap.New()
	s.Problems[id] = prob
	s.MaxProblemId = max(s.MaxProblemId, id)
	return prob
}

func (s *ServeCFG) Path() string {
	return filepath.Join(config.GlobalConfig.JudgesDir, fmt.Sprintf("%06d", s.ContestID))
}

func (s *ServeCFG) ProblemPath(id int) (string, error) {
	prob, ok := s.Problems[id]
	if !ok {
		return "", fmt.Errorf("no problem id %d in contest %d", id, s.ContestID)
	}
	probVal, ok := prob.Get("internal_name")
	if !ok {
		return "", fmt.Errorf("dir is not specified for problem id %d in contest %d", id, s.ContestID)
	}
	probStr, ok := probVal.(string)
	if !ok {
		return "", fmt.Errorf("unknown probStr for problem id %d in contest %d", id, s.ContestID)
	}

	return filepath.Join(s.Path(), "Problems", probStr), nil
}

func (s *ServeCFG) Write() error {
	b := &bytes.Buffer{}

	s.Global.Write(b)

	for _, lang := range s.languages {
		b.WriteString("[language]\n")
		lang.Write(b)
	}

	for _, prob := range s.abstractProblems {
		b.WriteString("[problem]\n")
		prob.Write(b)
	}

	for id := 0; id <= s.MaxProblemId; id++ {
		prob, ok := s.Problems[id]
		if !ok {
			continue
		}
		b.WriteString("[problem]\n")
		prob.Write(b)
	}

	for _, tester := range s.testers {
		b.WriteString("[tester]\n")
		tester.Write(b)
	}

	for conf, items := range s.other {
		for _, item := range items {
			b.WriteString(fmt.Sprintf("[%s]\n", conf))
			item.Write(b)
		}
	}

	return os.WriteFile(filepath.Join(s.Path(), "conf", "serve.cfg"), b.Bytes(), 0664)
}

var GENERIC = "Generic"

func (s *ServeCFG) addSection(currSectionName string, currSection *orderedmap.OrderedMap) {
	switch currSectionName {
	case "Global":
		s.Global = currSection
	case "language":
		s.languages = append(s.languages, currSection)
	case "tester":
		s.testers = append(s.testers, currSection)
	case "problem":
		if currSection.GetBool("abstract") {
			s.abstractProblems = append(s.abstractProblems, currSection)
			if currSection.GetStr("short_name") == GENERIC {
				s.HasGeneric = true
			}
		} else {
			s.Problems[currSection.GetInt("id")] = currSection
			s.MaxProblemId = max(s.MaxProblemId, currSection.GetInt("id"))
		}
	default:
		s.other[currSectionName] = append(s.other[currSectionName], currSection)
	}
}

func (s *ServeCFG) Clone() *ServeCFG {
	clonned := &ServeCFG{
		Global:    s.Global.Clone(),
		languages: s.languages, // We do not modify them
		testers:   s.testers,   // We do not modify them
		other:     s.other,     // We do not modify it

		Problems:         make(map[int]*orderedmap.OrderedMap),
		abstractProblems: s.abstractProblems, // We do not modify it

		ContestID: s.ContestID,

		HasGeneric:   s.HasGeneric,
		MaxProblemId: s.MaxProblemId,
	}

	for probId, prob := range s.Problems {
		clonned.Problems[probId] = prob.Clone()
	}

	return clonned
}
