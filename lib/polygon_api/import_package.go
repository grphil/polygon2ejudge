package polygon_api

import (
	"fmt"
)

func ImportPackage(problemUrl string, outPath string) error {
	r := newPostRequest()
	r.SetQueryParams(map[string]string{
		"type": "linux",
	})
	r.SetOutput(outPath)

	resp, err := r.Post(problemUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("exit code %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return nil
}
