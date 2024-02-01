package polygon_api

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

func ImportContest(constedID int, dir string) (*PContestXML, error) {
	xmlPath := filepath.Join(dir, "contest.xml")

	r := newPostRequest()
	r.SetOutput(xmlPath)
	resp, err := r.Post(fmt.Sprintf("https://polygon.codeforces.com/c/%d/contest.xml", constedID))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("exit code %d, body: %s", resp.StatusCode(), resp.Body())
	}

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, err
	}

	contest := &PContestXML{}
	err = xml.Unmarshal(xmlData, contest)
	if err != nil {
		return nil, err
	}

	// TODO: Add contest statements download
	return contest, nil
}
