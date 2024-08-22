package polygon_api

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"polygon2ejudge/lib/config"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func fixApiValues(method string, values url.Values) url.Values {
	tm := time.Now().Unix()
	values.Set("time", strconv.FormatInt(tm, 10))
	values.Set("apiKey", config.UserConfig.ApiKey)

	builder := bytes.Buffer{}
	rand := "000000"
	builder.WriteString(rand)
	builder.WriteRune('/')
	builder.WriteString(method)
	builder.WriteRune('?')
	builder.WriteString(values.Encode())
	builder.WriteRune('#')
	builder.WriteString(config.UserConfig.ApiSecret)

	endoded := sha512.Sum512(builder.Bytes())

	values.Set("apiSig", rand+hex.EncodeToString(endoded[:]))
	return values
}

const kRetryCount = 10

func newPostRequest() *resty.Request {
	c := resty.New()
	c.SetFormData(map[string]string{
		"login":    config.UserConfig.PolygonLogin,
		"password": config.UserConfig.PolygonPassword,
	})
	c.SetRetryCount(kRetryCount)
	c.SetRetryMaxWaitTime(time.Second * time.Duration(129))
	return c.R()
}
