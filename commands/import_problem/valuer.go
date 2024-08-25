package import_problem

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/config"
	"polygon2ejudge/lib/orderedmap"
	"slices"
	"strconv"
	"strings"
)

func (t *ImportTask) buildValuer() error {
	acm := t.ServeCFG.Global.GetStr("score_system") == "acm"
	err := t.processTests(acm)
	if err != nil {
		return fmt.Errorf("%s error while processing tests %s", t.InternalName, err.Error())
	}

	valuerOptions := orderedmap.New()

	err = t.processGroups(acm, valuerOptions)
	if err != nil {
		return err
	}

	fmt.Println("Processed tests and groups")

	return nil
}

func (t *ImportTask) processGroups(
	acm bool,
	valuerOptions *orderedmap.OrderedMap,
) error {
	userScore := 0
	fullScore := 0
	var openTests []string
	var finalOpenTests []string
	var testScoreList []string

	valuerCFG := &bytes.Buffer{}
	valuerCFG.WriteString("global {\n    stat_to_judges 1;\n    stat_to_users 0;\n}\n\n")

	for _, g := range t.groups {
		if len(g.tests) == 0 {
			return fmt.Errorf("group %s has no tests", g.name)
		}
		tests := fmt.Sprintf("%d-%d", g.tests[0], g.tests[len(g.tests)-1])

		valuerCFG.WriteString(fmt.Sprintf("group %s {\n", g.name))
		valuerCFG.WriteString(fmt.Sprintf("    tests %s;\n", tests))

		var testScores []string

		switch g.scoreType {
		case ScoreTypeLowestScore:
			valuerCFG.WriteString("    use_lowest_test_score;\n")
			fallthrough
		case ScoreTypeEachTest:
			valuerCFG.WriteString(fmt.Sprintf("    test_score %d;\n", g.score))
			for _ = range g.tests {
				testScores = append(testScores, strconv.Itoa(g.score))
			}
		case ScoreTypeComplete:
			valuerCFG.WriteString(fmt.Sprintf("    score %d;\n", g.score))
			for _ = range g.tests {
				testScores = append(testScores, "0")
			}
			testScores[len(testScores)-1] = strconv.Itoa(g.score)
		}
		testScoreList = append(testScoreList, strings.Join(testScores, " "))

		score := g.score
		if g.scoreType == ScoreTypeEachTest {
			score *= len(g.tests)
		}

		if g.feedbackType != FeedbackTypeHidden {
			userScore += score
		}
		fullScore += score

		if len(g.requires) > 0 {
			valuerCFG.WriteString(fmt.Sprintf(
				"    requires %s;\n",
				strings.Join(g.requires, ", "),
			))
		}

		if g.setsMarked {
			g.requires = append(g.requires, g.name)
			valuerCFG.WriteString(fmt.Sprintf(
				"    sets_marked_if_passed %s;\n",
				strings.Join(g.requires, ", "),
			))
		}

		var openTestsMode string
		switch g.feedbackType {
		case FeedbackTypeHidden:
			openTestsMode = "hidden"
			valuerCFG.WriteString("    offline;\n")
		case FeedbackTypeExists:
			openTestsMode = "exists"
		case FeedbackTypeBrief:
			openTestsMode = "brief"
		case FeedbackTypeFull:
			openTestsMode = "full"
		}
		openTests = append(openTests, fmt.Sprintf(
			"%d-%d:%s",
			g.tests[0],
			g.tests[len(g.tests)-1],
			openTestsMode,
		))

		finalOpenTests = append(finalOpenTests, fmt.Sprintf(
			"%d-%d:%s",
			g.tests[0],
			g.tests[len(g.tests)-1],
			"full",
		))

		valuerCFG.WriteString("}\n\n")
	}

	valuerOptions.Set("open_tests", strings.Join(openTests, ", "))
	valuerOptions.Set("final_open_tests", strings.Join(finalOpenTests, ", "))
	if !acm {
		t.ServeCFG.Global.Set("separate_user_score", true)

		valuerOptions.Set("full_score", fullScore)
		valuerOptions.Set("full_user_score", userScore)
		valuerOptions.Set("test_score_list", strings.Join(testScoreList, "  "))
		valuerOptions.Set("valuer_cmd", "../gvaluer")
		valuerOptions.Set("interactive_valuer", true)
		valuerOptions.Set("valuer_sets_marked", true)
		valuerOptions.Set("olympiad_mode", true)
		valuerOptions.Set("run_penalty", 0)

		valuerR, err := os.Open(config.GlobalConfig.GvaluerPath)
		if err != nil {
			return err
		}
		defer valuerR.Close()

		valuerPath := filepath.Join(t.tmpDir, "gvaluer")
		valuerW, err := os.OpenFile(valuerPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
		if err != nil {
			return err
		}

		_, err = io.Copy(valuerW, valuerR)
		if err != nil {
			return err
		}
		valuerW.Close()

		err = t.Transaction.MovePath(valuerPath, filepath.Join(t.ServeCFG.Path(), "problems", "gvaluer"))
		if err != nil {
			return err
		}
	}

	t.config.Update(valuerOptions)

	valuerOptions.WritePrefix(valuerCFG, "# ")
	err := os.WriteFile(
		filepath.Join(t.ProbDir, "valuer.cfg"),
		valuerCFG.Bytes(),
		0664,
	)
	if err != nil {
		return fmt.Errorf("error while writing valuer.cfg, error: %s", err.Error())
	}
	return nil
}

func (t *ImportTask) processTests(acm bool) error {
	hasGroups := len(t.testset.Groups.Groups) > 0
	for i, test := range t.testset.Tests.Tests {
		if len(test.Group) == 0 {
			if hasGroups {
				return fmt.Errorf("test %d has no group", i+1)
			}
			if test.Sample {
				test.Group = "0"
			} else {
				test.Group = "1"
			}
			test.Points = 0.0
		}

		if acm {
			if len(t.groups) == 0 || t.LastGroup().samples != test.Sample {
				err := t.AddGroupInfo(test, nil, acm)
				if err != nil {
					return err
				}
			}
		} else if len(t.groups) == 0 || t.LastGroup().name != test.Group {
			if slices.ContainsFunc(t.groups, func(g *GroupInfo) bool {
				return g.name == test.Group
			}) {
				return fmt.Errorf("wrong tests for group %s, tests are not consequent", test.Group)
			}

			groupId := slices.IndexFunc(t.testset.Groups.Groups, func(group *XGroup) bool {
				return test.Group == group.Name
			})

			var group *XGroup
			if groupId != -1 {
				group = t.testset.Groups.Groups[groupId]
			} else if hasGroups {
				return fmt.Errorf("test %d belongs to group %s which is not found", i+1, test.Group)
			} else {
				group = &XGroup{
					Name:           test.Group,
					FeedbackPolicy: "icpc",
					PointsPolicy:   "complete-group",
					Points:         0.0,
				}

				if !test.Sample {
					group.Dependencies.Dependencies = append(group.Dependencies.Dependencies, &XGroupDependency{Group: "0"})
				}
			}

			err := t.AddGroupInfo(test, group, acm)
			if err != nil {
				return err
			}
		}

		err := t.LastGroup().AddTest(i, test)
		if err != nil {
			return err
		}
	}
	if !acm && !hasGroups && len(t.groups) > 0 {
		t.LastGroup().score = 100
		t.markLastGroup()
	}
	return nil
}
