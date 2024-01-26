package import_problem

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (t *ImportTask) openZip() error {

	return nil
}

func (t *ImportTask) extractProblemFiles() error {
	archive, err := zip.OpenReader(t.packagePath)
	if err != nil {
		return err
	}

	t.archive = archive

	t.addDeferFunc(func() {
		t.archive.Close()
	})

	t.probDir = filepath.Join(t.tmpDir, "problem")
	err = os.MkdirAll(filepath.Dir(t.probDir), 0774)
	if err != nil {
		return err
	}

	err = t.moveFileName("problem.xml", "")
	if err != nil {
		return err
	}

	probXMLData, err := os.ReadFile(filepath.Join(t.probDir, "problem.xml"))
	if err != nil {
		return err
	}

	t.problemXML = &XProblemXML{}

	err = xml.Unmarshal(probXMLData, t.problemXML)
	if err != nil {
		return fmt.Errorf("can not parse problem xml, error: %s", err.Error())
	}

	fmt.Printf("Parsed problem.xml for problem %s id %d\n", t.problemXML.ShortName, *t.PolygonProbId)

	err = t.extractSolutions()
	if err != nil {
		return fmt.Errorf("error while extracting solutions, error: %s", err.Error())
	}

	err = t.extractAssets()
	if err != nil {
		return err
	}

	err = t.extractAllFiles("tests", "tests")
	if err != nil {
		return fmt.Errorf("can not extract tests, error: %s", err.Error())
	}

	err = t.extractAllFiles("documents", "documents")
	if err != nil {
		return fmt.Errorf("can not extract documents, error: %s", err.Error())
	}

	fmt.Println("Extracted all problem files from zip archive")
	return nil
}

func (t *ImportTask) extractSolutions() error {
	solutionsPath := "solutions1"
	if *t.CompileAllSolutions {
		solutionsPath = "solutions"
	}

	for _, solution := range t.problemXML.Assets.Solutions.Solutions {
		err := t.moveFileName(solution.Source.Path, solutionsPath)
		if err != nil {
			return err
		}

		if solution.Tag == "main" {
			name := filepath.Base(solution.Source.Path)
			err = copyFile(
				filepath.Join(t.probDir, solutionsPath, name),
				filepath.Join(t.probDir, name),
			)
		}

	}
	return nil
}

func (t *ImportTask) extractAssets() error {
	err := t.moveFileName(t.problemXML.Assets.Checker.Source.Path, "")
	if err != nil {
		return fmt.Errorf("can not extract checker, error: %s", err.Error())
	}

	if t.problemXML.Assets.Interactor != nil {
		err = t.moveFileName(t.problemXML.Assets.Interactor.Source.Path, "")
		if err != nil {
			return fmt.Errorf("can not extract interactor, error: %s", err.Error())
		}
	}

	for _, resource := range t.problemXML.Files.Resources.Resources {
		err = t.moveFileName(resource.Path, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ImportTask) extractAllFiles(prefix string, dst string) error {
	for _, file := range t.archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name, prefix) {
			dstPath, err := filepath.Rel(prefix, file.Name)
			if err != nil {
				panic(err) // If Name has prefix, this is impossible, it is impossible
			}

			dstPath = filepath.Join(t.probDir, dst, dstPath)
			err = os.MkdirAll(filepath.Dir(dstPath), 0774)
			if err != nil {
				return err
			}

			err = moveSingleFile(file, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *ImportTask) moveFileName(srcPath string, dstDir string) error {
	return t.moveFile(srcPath, dstDir, filepath.Base(srcPath))
}

func (t *ImportTask) moveFile(srcPath string, dstDir string, name string) error {
	dstDir = filepath.Join(t.probDir, dstDir)
	err := os.MkdirAll(dstDir, 0774)
	if err != nil {
		return fmt.Errorf("can not create dir %s for file %s, error: %s", dstDir, srcPath, err.Error())
	}

	for _, file := range t.archive.File {
		if file.Name == srcPath {
			err = moveSingleFile(file, filepath.Join(dstDir, name))
			if err != nil {
				return fmt.Errorf("can not extract zip file %s to dir %s, error: %s", srcPath, dstDir, err.Error())
			}
			return nil
		}
	}
	return fmt.Errorf("file with path %s not found", srcPath)
}

func moveSingleFile(src *zip.File, dst string) error {
	w, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return err
	}
	defer w.Close()

	r, err := src.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}
