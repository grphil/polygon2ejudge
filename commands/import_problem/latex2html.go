package import_problem

import (
	"fmt"
	"github.com/anaskhan96/soup"
	html2 "golang.org/x/net/html"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func (t *ImportTask) processTex(f string) string {
	html, err := t.latex2HTML(f)
	if err != nil {
		fmt.Printf("Warning: can not convert %s to html, error: %s\n", f, err.Error())
		return ""
	}

	html, err = t.fixHTML(html)
	if err != nil {
		fmt.Printf("Warning: can not process %s html, error: %s\n", f, err.Error())
		return ""
	}
	return html
}

func (t *ImportTask) latex2HTML(f string) (string, error) {
	tmpPath := filepath.Join(t.tmpDir, "converter")
	err := os.MkdirAll(tmpPath, 0774)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpPath)

	contentB, err := os.ReadFile(filepath.Join(t.statementPath, f))
	if err != nil {
		return "", err
	}

	content := t.fixLatex(string(contentB))

	texPath := filepath.Join(tmpPath, "a.tex")
	err = os.WriteFile(texPath, []byte(content), 0664)
	if err != nil {
		return "", err
	}
	resPath := filepath.Join(tmpPath, "a.html")

	cmd := exec.Command("pandoc", "-f", "latex", "-t", "html", "--mathjax", texPath, "-o", resPath)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("can not convert %s from latex to html\nerror: %s\n%s\n", f, err.Error(), stdoutStderr)
	}

	html, err := os.ReadFile(resPath)
	if err != nil {
		return "", err
	}
	return string(html), nil
}

func (t *ImportTask) fixLatex(content string) string {
	for strings.Contains(content, "\\input") {
		content = t.processTexInput(content)
		fmt.Println(content)
	}

	content = strings.ReplaceAll(content, "\\t{", "\\texttt{")
	content = strings.ReplaceAll(content, "<<", "«")
	content = strings.ReplaceAll(content, ">>", "»")

	// fuck latex tables
	content = strings.ReplaceAll(content, "\\parbox", "")

	runes := []rune(content)
	c := strings.Builder{}

	cmp := func(i int, s string) bool {
		l := len([]rune(s))
		ss := string(runes[i:min(len(runes), i+l)])
		return ss == s
	}

	inTable := false
	balance := 0
	for i := 0; i < len(runes); i++ {
		if !inTable {
			c.WriteRune(runes[i])
			if cmp(i, "\\begin{tabular}") {
				inTable = true
				balance = 0
			}
			continue
		} else {
			if cmp(i, "\\{") {
				c.WriteString("\\{")
				i++
				continue
			}
			if cmp(i, "\\}") {
				c.WriteString("\\}")
				i++
				continue
			}
			if balance > 0 && cmp(i, "\\\\") {
				c.WriteRune(' ')
				i++
				continue
			}
			c.WriteRune(runes[i])
			if cmp(i, "\\end{tabular}") {
				inTable = false
			}
			if runes[i] == '{' {
				balance++
			}
			if runes[i] == '}' {
				balance--
			}
		}
	}

	return c.String()
}

func (t *ImportTask) processTexInput(content string) string {
	inputPos := strings.Index(content, "\\input")
	contentPrefix := content[:inputPos]
	contentSuffix := content[inputPos+6:]

	if len(contentSuffix) == 0 || contentSuffix[0] != '{' {
		return contentPrefix + contentSuffix
	}

	fileEndIndex := strings.Index(contentSuffix, "}")
	if fileEndIndex == -1 {
		return contentPrefix + contentSuffix
	}

	fileName := contentSuffix[1:fileEndIndex]
	contentSuffix = contentSuffix[fileEndIndex+1:]

	if !strings.HasSuffix(fileName, ".tex") {
		fileName += ".tex"
	}

	inputFileContentB, err := os.ReadFile(filepath.Join(t.statementPath, fileName))
	if err != nil {
		return contentPrefix + contentSuffix
	}

	return contentPrefix + "\n\n" + string(inputFileContentB) + "\n\n" + contentSuffix
}

const kStyle = "width: auto; max-width: max(50%, 400px); height: auto; max-height: 100%;"

func (t *ImportTask) fixHTML(s string) (string, error) {
	r := regexp.MustCompile(`\[[0-9]+(cm|mm)\]`)
	s = r.ReplaceAllString(s, "")

	r = regexp.MustCompile(`<span>[0-9]+(cm|mm)</span>`)
	s = r.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, "<table>", "<table class=\"statements\">")

	html := soup.HTMLParse(s)
	if html.Error != nil {
		return "", html.Error
	}

	images := html.FindAll("img")
	for _, img := range images {
		for i := range img.Pointer.Attr {
			if img.Pointer.Attr[i].Key != "src" {
				continue
			}
			imgName := img.Pointer.Attr[i].Val
			err := copyFile(
				filepath.Join(t.statementPath, imgName),
				filepath.Join(t.ProbDir, "attachments", imgName),
			)
			if err != nil {
				fmt.Printf("Warning: can not process image %s, error: %s", imgName, err.Error())
			} else {
				img.Pointer.Attr[i].Val = "${getfile}=" + imgName
			}
			break
		}
		img.Pointer.Attr = append(img.Pointer.Attr, html2.Attribute{
			Key: "style",
			Val: kStyle,
		})
	}

	embeds := html.FindAll("embed")
	for _, embed := range embeds {
		embed.Pointer.Data = "img"
		for i := range embed.Pointer.Attr {
			if embed.Pointer.Attr[i].Key != "src" {
				continue
			}
			epsName := embed.Pointer.Attr[i].Val
			imgName, err := t.convertEps(epsName)
			if err != nil {
				fmt.Printf("Warning: can not convert image %s, error: %s", epsName, err.Error())
			} else {
				embed.Pointer.Attr[i].Val = "${getfile}=" + imgName
			}
			break
		}
		embed.Pointer.Attr = append(embed.Pointer.Attr, html2.Attribute{
			Key: "style",
			Val: kStyle,
		})
	}

	s = html.HTML()
	s, _ = strings.CutPrefix(s, "<html><head></head><body>")
	s, _ = strings.CutSuffix(s, "</body></html>")
	return s, nil
}

func (t *ImportTask) convertEps(name string) (string, error) {
	tmpPath := filepath.Join(t.tmpDir, "converter")
	err := os.MkdirAll(tmpPath, 0774)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpPath)

	epsPath := filepath.Join(tmpPath, name)
	err = copyFile(filepath.Join(t.statementPath, name), epsPath)
	if err != nil {
		return "", err
	}

	resName := fmt.Sprintf("%d.png", rand.Intn(1000000000))
	resTmpPath := filepath.Join(tmpPath, resName)

	cmd := exec.Command(
		"gs",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-dEPSCrop",
		"-r600",
		"-sDEVICE=pngalpha",
		fmt.Sprintf("-sOutputFile=%s", resTmpPath),
		epsPath,
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("can not convert image %s to png \nerror: %s\n%s\n", name, err.Error(), stdoutStderr)
	}

	err = copyFile(resTmpPath, filepath.Join(t.statementPath, "attachments", resName))
	if err != nil {
		return "", err
	}
	return resName, nil
}

func copyFile(src string, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0774)
	if err != nil {
		return err
	}

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}
