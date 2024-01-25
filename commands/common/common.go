package common

import "github.com/hellflame/argparse"

type TaskCommon struct {
	ContestId        *int
	ResetCredentials *bool
	NoOffline        *bool
	Abstract         *string
}

func (t *TaskCommon) AddCommonOptions(parser *argparse.Parser, hasImport bool) {
	t.ContestId = parser.Int("", "contest_id", &argparse.Option{
		Help:       "Id of ejudge contest to add problem",
		Required:   true,
		Positional: true,
	})
	t.ResetCredentials = parser.Flag("", "reset-credentials", &argparse.Option{
		Help: "Reset polygon api and ejudge login data",
	})
	if hasImport {
		t.Abstract = parser.String("", "abstract", &argparse.Option{
			Help: "Abstract problem to import all configs from",
		})
		t.NoOffline = parser.Flag("", "no-offline", &argparse.Option{
			Help: "No offline groups for task",
		})
	}

}
