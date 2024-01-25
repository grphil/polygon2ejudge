package import_problem

import (
	"fmt"
	"polygon2ejudge/lib/config"
	"slices"
)

type ScoreType int

type FeedbackType int

const (
	ScoreTypeComplete ScoreType = iota
	ScoreTypeEachTest
	ScoreTypeLowestScore
)

const (
	FeedbackTypeFull FeedbackType = iota
	FeedbackTypeBrief
	FeedbackTypeExists
	FeedbackTypeHidden
)

type GroupInfo struct {
	name     string
	requires []string

	tests []int

	scoreType    ScoreType
	feedbackType FeedbackType

	samples    bool
	score      int
	acm        bool
	setsMarked bool
}

func (g *GroupInfo) AddTest(i int, test *XTest) error {
	if len(g.tests) == 0 {
		g.tests = append(g.tests, i+1)
	} else {
		if g.tests[len(g.tests)-1] != i {
			return fmt.Errorf("wrong tests for group %s, tests are not consequent", g.name)
		}
		g.tests = append(g.tests, i+1)
	}

	if g.samples != test.Sample {
		return fmt.Errorf("group %s contains both samples and normal tests", g.name)
	}

	if g.acm {
		return nil
	}

	if g.name != test.Group {
		return fmt.Errorf("adding test with group %s to group %s", test.Group, g.name)
	}

	switch g.scoreType {
	case ScoreTypeEachTest, ScoreTypeLowestScore:
		if int(test.Points) != g.score {
			return fmt.Errorf("group %s has tests with different scores", g.name)
		}
	case ScoreTypeComplete:
		g.score += int(test.Points)
	}
	return nil
}

func (t *ImportTask) LastGroup() *GroupInfo {
	return t.groups[len(t.groups)-1]
}

func (t *ImportTask) AddGroupInfo(test *XTest, group *XGroup, acm bool) error {
	g := &GroupInfo{
		name:    test.Group,
		samples: test.Sample,
		acm:     acm,
	}

	if acm {
		if g.samples {
			g.feedbackType = FeedbackTypeFull
		} else {
			g.feedbackType = FeedbackTypeBrief
		}
		t.groups = append(t.groups, g)
		return nil
	}

	if test.Group != group.Name {
		return fmt.Errorf("adding test with group %s to group %s", test.Group, group.Name)
	}

	for _, dep := range group.Dependencies.Dependencies {
		if !slices.ContainsFunc(t.groups, func(info *GroupInfo) bool {
			return info.name == dep.Group
		}) {
			return fmt.Errorf("group %s depends on group %s, which occurs later", test.Group, dep.Group)
		}
		g.requires = append(g.requires, dep.Group)
	}

	switch group.PointsPolicy {
	case "complete-group":
		g.scoreType = ScoreTypeComplete
		g.score = 0
	case "each-test":
		g.scoreType = ScoreTypeEachTest
		g.score = int(test.Points)
	case "lowest-score":
		g.scoreType = ScoreTypeLowestScore
		g.score = int(test.Points)
	default:
		return fmt.Errorf("unsupported group points polycy %s", group.PointsPolicy)
	}

	if g.samples {
		g.feedbackType = FeedbackTypeFull
	} else {
		switch group.FeedbackPolicy {
		case "none":
			if *t.NoOffline {
				g.feedbackType = FeedbackTypeExists
			} else {
				g.feedbackType = FeedbackTypeHidden

				t.markLastGroup()
			}
		case "points":
			g.feedbackType = FeedbackTypeExists
		case "icpc":
			g.feedbackType = FeedbackTypeBrief
		case "complete":
			if config.FULL_REPORT_ONLY_SAMPLES {
				g.feedbackType = FeedbackTypeBrief
			} else {
				g.feedbackType = FeedbackTypeFull
			}
		default:
			return fmt.Errorf("unsupported feedback policy %s", group.FeedbackPolicy)
		}
	}

	t.groups = append(t.groups, g)

	return nil
}

func (t *ImportTask) markLastGroup() {
	if len(t.groups) == 0 {
		return
	}

	g := t.LastGroup()
	if g.feedbackType == FeedbackTypeHidden {
		return
	}

	for _, info := range t.groups[:len(t.groups)-1] {
		if info.feedbackType == FeedbackTypeHidden {
			fmt.Printf("Warning: offline group %s before online group %s\n", g.name, info.name)
			return
		}

		if !slices.ContainsFunc(g.requires, func(name string) bool {
			return info.name == name
		}) {
			fmt.Printf("Warning: last online group %s does not depend on group %s\n", info.name, g.name)
			return
		}
	}

	g.setsMarked = true
}
