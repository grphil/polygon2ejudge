package polygon_api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"strconv"
)

const PACKAGES_METHOD = "problem.packages"
const PACKAGE_METHOD = "problem.package"

func (p *PolygonApi) ImportPackage(probIDInt int, outPath string) error {
	probID := strconv.Itoa(probIDInt)
	packageID, err := p.getPackages(probID)
	if err != nil {
		return fmt.Errorf("%s: %s", PACKAGES_METHOD, err.Error())
	}

	err = p.downloadPackage(probID, strconv.Itoa(packageID), outPath)
	if err != nil {
		return fmt.Errorf("%s: %s", PACKAGE_METHOD, err.Error())
	}
	fmt.Printf("Downloaded problem %d package to %s\n", probIDInt, outPath)
	return nil
}

func (p *PolygonApi) downloadPackage(probID string, packageID string, outPath string) error {
	values := url.Values{}
	values.Set("problemId", probID)
	values.Set("packageId", packageID)
	values.Set("type", "linux")
	values = p.fixValues(PACKAGE_METHOD, values)

	client := resty.New()
	r := client.R()
	r.SetQueryParamsFromValues(values)
	r.SetOutput(outPath)
	resp, err := r.Get(POLYGON_URL + PACKAGE_METHOD)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("exit code %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return nil
}

type PPackagesList struct {
	Status  string            `json:"status"`
	Comment string            `json:"comment"`
	Result  []*PSinglePackage `json:"result"`
}

type PSinglePackage struct {
	ID       int    `json:"id"`
	Revision int    `json:"revision"`
	Type     string `json:"type"`
}

func (p *PolygonApi) getPackages(probID string) (int, error) {
	values := url.Values{}
	values.Set("problemId", probID)
	values = p.fixValues(PACKAGES_METHOD, values)

	client := resty.New()
	r := client.R()
	r.SetQueryParamsFromValues(values)
	res, err := r.Get(POLYGON_URL + PACKAGES_METHOD)
	if err != nil {
		return 0, err
	}
	if res.StatusCode() != 200 {
		return 0, fmt.Errorf("exit code %d, body: %s", res.StatusCode(), res.String())
	}

	packages := &PPackagesList{}
	err = json.Unmarshal(res.Body(), packages)
	if err != nil {
		return 0, fmt.Errorf("can not parse json response, error: %s", err.Error())
	}

	if packages.Status != "OK" {
		return 0, fmt.Errorf("responded with %s, comment: %s", packages.Status, packages.Comment)
	}

	var bestPackage *PSinglePackage
	for _, pack := range packages.Result {
		if pack.Type != "linux" {
			continue
		}

		if bestPackage == nil || pack.Revision > bestPackage.Revision {
			bestPackage = pack
		}
	}

	if bestPackage == nil {
		return 0, fmt.Errorf("no package created for problem %s", probID)
	}

	return bestPackage.ID, nil
}
