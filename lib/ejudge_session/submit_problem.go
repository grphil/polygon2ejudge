package ejudge_session

import (
	"strconv"
	"strings"
)

func (s *TEjudgeSession) SubmitProblem(
	data string,
	problemID int,
	langID int,
) error {
	r := s.client.R()
	r.SetMultipartFormData(map[string]string{
		"SID":       s.sid,
		"problem":   strconv.Itoa(problemID),
		"eoln_type": "1",
		"lang_id":   strconv.Itoa(langID),
		"action_40": "Send!",
	})

	reader := strings.NewReader(data)
	r.SetFileReader("file", "file", reader)
	_, err := r.Post(newJudgeUrl())
	if err != nil {
		return err
	}
	return nil
}
