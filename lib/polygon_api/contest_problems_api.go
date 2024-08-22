package polygon_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"polygon2ejudge/lib/config"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const CONTEST_METHOD = "contest.problems"

type PSingleProblem struct {
	ID            int    `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Deleted       bool   `json:"deleted"`
	Favourite     bool   `json:"favourite"`
	AccessType    string `json:"accessType"`
	Revision      int    `json:"revision"`
	LatestPackage int    `json:"latestPackage"`
	Modified      bool   `json:"modified"`
}

type PProblemsList struct {
	Status  string                     `json:"status"`
	Comment string                     `json:"comment"`
	Result  map[string]*PSingleProblem `json:"result"`
}

// Problem idx (A, B, C) -> Problem.
func GetProblemsInContest(contestID int) (map[string]*PSingleProblem, error) {
	values := url.Values{}
	values.Set("contestId", strconv.Itoa(contestID))
	values = fixApiValues(CONTEST_METHOD, values)

	client := resty.New()
	r := client.R()
	r.SetQueryParamsFromValues(values)
	res, err := r.Get(config.GlobalConfig.PolygonUrl + CONTEST_METHOD)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("exit code %d, body: %s", res.StatusCode(), res.String())
	}

	var plist PProblemsList
	if err := json.Unmarshal(res.Body(), &plist); err != nil {
		return nil, err
	}
	if plist.Status != "OK" {
		return nil, fmt.Errorf("responded with %s, comment: %s", plist.Status, plist.Comment)
	}

	return plist.Result, nil
}
