package ejudge_session

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http/cookiejar"
	"polygon2ejudge/lib/config"
	"strconv"
)

type TEjudgeSession struct {
	client    *resty.Client
	sid       string
	contestID int
}

func NewEjudgeSession(contestID int) (*TEjudgeSession, error) {
	client := resty.New()
	jar, _ := cookiejar.New(nil)
	client.SetCookieJar(jar)

	r := client.R()
	r.SetFormData(map[string]string{
		"login":      config.UserConfig.EjudgeLogin,
		"password":   config.UserConfig.EjudgePassword,
		"contest_id": strconv.Itoa(contestID),
		"role":       "1",
		"language":   "0",
		"action_2":   "",
	})

	resp, err := r.Post(newJudgeUrl())
	if err != nil {
		return nil, err
	}

	u := resp.RawResponse.Request.URL.Query()
	sid := u.Get("SID")
	if len(sid) == 0 {
		return nil, errors.New("not authorised, no sid")
	}

	ejudgeSession := &TEjudgeSession{
		client:    client,
		sid:       sid,
		contestID: contestID,
	}
	fmt.Printf("Connected to ejudge contest %d\n", contestID)
	return ejudgeSession, nil
}

func newJudgeUrl() string {
	return config.GlobalConfig.EjudgeUrl + "/cgi-bin/new-judge"
}
