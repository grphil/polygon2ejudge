package import_problem

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"polygon2ejudge/lib/orderedmap"
	"polygon2ejudge/lib/serve_cfg"
	"strconv"
	"strings"
)

func (t *ImportTask) fillInConfig() error {
	if *t.EjudgeId == -1 {
		id := t.serveCFG.MaxProblemId + 1
		t.EjudgeId = &id
	}
	t.config = t.serveCFG.Problem(*t.EjudgeId)
	t.problemOnlyConfig = orderedmap.New()
	t.addDeferFunc(t.exportConfig)

	t.config.Set("id", *t.EjudgeId)

	abstract := t.Abstract
	if len(*abstract) == 0 && !*t.NoGenericParent && t.serveCFG.HasGeneric {
		abstract = &serve_cfg.GENERIC
	}
	if len(*abstract) > 0 {
		t.config.Set("super", *t.Abstract)
	}

	t.setNames()

	t.config.Set("extid", fmt.Sprintf("polygon:%d", *t.PolygonProbId))
	t.problemOnlyConfig.Set("revision", t.problemXML.Revision)

	t.config.Set("use_stdin", true)
	t.config.Set("use_stdout", true)
	if len(t.problemXML.Judging.InputFile) > 0 {
		t.config.Set("input_file", t.problemXML.Judging.InputFile)
		t.config.Set("combined_stdin", true)
	}
	if len(t.problemXML.Judging.OutputFile) > 0 {
		t.config.Set("output_file", t.problemXML.Judging.OutputFile)
		t.config.Set("combined_stdout", true)
	}

	t.config.Set("test_pat", "%02d")
	t.config.Set("use_corr", true)
	t.config.Set("corr_pat", "%02d.a")

	err := t.setLimits()
	if err != nil {
		return err
	}

	t.config.Set("check_cmd", getFileName(t.problemXML.Assets.Checker.Source))
	if t.problemXML.Assets.Interactor != nil {
		t.config.Set("interactor_cmd", getFileName(t.problemXML.Assets.Interactor.Source))
	}

	if !*t.NoCompileMainSolution {
		for _, s := range t.problemXML.Assets.Solutions.Solutions {
			if s.Tag == "main" {
				t.config.Set("solution_cmd", getFileName(s.Source))
			}
		}
	}

	t.config.Set("enable_testlib_mode", true)

	if !*t.NoTextareaInput {
		t.config.Set("enable_text_form", true)
	}

	if !*t.NoCustomRun {
		t.config.Set("enable_user_input", true)
	}

	t.setCustomOptions()

	return nil
}

func (t *ImportTask) setCustomOptions() {
	file, err := os.Open(filepath.Join(t.probDir, "documents/description.txt"))
	if err != nil {
		return
	}
	fmt.Println("description.txt file provided, will try read configs from it")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())

		idx := strings.Index(txt, " ")
		if idx == -1 {
			continue
		}
		var first, second, third string

		first = strings.TrimSpace(txt[:idx])
		txt = strings.TrimSpace(txt[idx:])

		idx = strings.Index(txt, " ")
		if idx == -1 {
			second = txt
		} else {
			second = strings.TrimSpace(txt[:idx])
			third = strings.TrimSpace(txt[idx:])
		}

		if first == "ejudge_config" {
			if len(third) == 0 {
				t.config.Set(second, orderedmap.NoVal(true))
			} else {
				t.config.Set(second, orderedmap.UnquotedStr(third))
			}
		}
		if first == "ejudge_remove_config" {
			t.config.Remove(second)
			t.problemOnlyConfig.Remove(second)
		}
		if first == "group_points" {
			for _, g := range t.testset.Groups.Groups {
				if g.Name == second {
					g.PointsPolicy = third
				}
			}
		}
	}
	if scanner.Err() != nil {
		fmt.Println("error found while reading description.txt ", err.Error())
	}
}

func (t *ImportTask) setLimits() error {
	for _, testset := range t.problemXML.Judging.Testsets {
		if testset.Name != "tests" {
			continue
		}

		timeLimit := testset.TimeLimit
		if timeLimit%1000 == 0 {
			t.config.Set("time_limit", timeLimit/1000)
		} else {
			t.config.Set("time_limit_millis", timeLimit)
		}
		t.config.Set("real_time_limit", max(5, (timeLimit*2+999)/1000))

		memoryLimit := testset.MemoryLimit / 1024
		var memoryLimitStr string

		if memoryLimit%1024 != 0 {
			memoryLimitStr = fmt.Sprintf("%dK", memoryLimit)
		} else {
			memoryLimit /= 1024
			if memoryLimit%1024 != 0 {
				memoryLimitStr = fmt.Sprintf("%dM", memoryLimit)
			} else {
				memoryLimitStr = fmt.Sprintf("%dG", memoryLimit/1024)
			}
		}
		t.config.Set("max_vm_size", orderedmap.UnquotedStr(memoryLimitStr))
		t.config.Set("max_stack_size", orderedmap.UnquotedStr(memoryLimitStr))
		t.testset = testset
		return nil
	}
	return fmt.Errorf("testset \"tests\" not found")
}

func (t *ImportTask) setNames() error {
	shortNames := make(map[string]bool)
	internalNames := make(map[string]bool)
	for _, prob := range t.serveCFG.Problems {
		shortNames[prob.GetStr("short_name")] = true
		internalNames[prob.GetStr("internal_name")] = true
	}

	if len(*t.ShortName) > 0 {
		_, ok := shortNames[*t.ShortName]
		if ok {
			t.ShortName = nil
		}
	} else {
		t.ShortName = nil
	}

	if t.ShortName == nil {
		for i := 'A'; i <= 'Z'; i++ {
			name := string(i)
			if _, ok := shortNames[name]; !ok {
				t.ShortName = &name
				break
			}
		}

		if t.ShortName == nil {
			for i := 1; true; i++ {
				name := strconv.Itoa(i)
				if _, ok := shortNames[name]; !ok {
					t.ShortName = &name
					break
				}
			}
		}
	}
	t.config.Set("short_name", *t.ShortName)

	internalName := t.problemXML.ShortName
	if _, ok := internalNames[internalName]; ok {
		for i := 2; true; i++ {
			name := fmt.Sprintf("%s-%d", internalName, i)
			if _, ok := internalNames[name]; !ok {
				internalName = name
				break
			}
		}
	}
	t.config.Set("internal_name", internalName)
	t.internalName = internalName
	if !*t.NoPackageSave {
		err := t.Transaction.MovePath(
			t.packagePath,
			filepath.Join(t.serveCFG.Path(), "download", fmt.Sprintf("%s.zip", t.internalName)),
		)
		if err != nil {
			return err
		}
	}

	var russianName string
	var englishName string
	var anyName string
	for _, name := range t.problemXML.Names.Names {
		if name.Language == "russian" {
			russianName = name.Value
		}
		if name.Language == "english" {
			englishName = name.Value
		}
		anyName = name.Value
	}
	if len(englishName) == 0 {
		englishName = russianName
	}
	if len(russianName) == 0 {
		russianName = englishName
	}
	if len(russianName) == 0 && len(englishName) == 0 {
		russianName = anyName
		englishName = anyName
	}
	t.config.Set("long_name", russianName)
	t.problemOnlyConfig.Set("long_name_en", englishName)
	return nil
}

func (t *ImportTask) exportConfig() {
	t.problemOnlyConfig.Update(t.config)

	b := &bytes.Buffer{}
	b.WriteString("# -*- coding: utf-8 -*-\\n\\n[problem]\n")
	t.problemOnlyConfig.Write(b)

	err := os.WriteFile(
		filepath.Join(t.probDir, "problem.cfg"),
		b.Bytes(),
		0664,
	)

	if err != nil {
		t.Transaction.SetError(fmt.Errorf("can not export problem config, error: %s", err.Error()))
	}
}

func getFileName(s XSource) string {
	file := filepath.Base(s.Path)
	dotPos := strings.Index(file, ".")
	if dotPos != -1 {
		return file[:dotPos]
	}
	return file
}
