package polygon_api

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"
)

func (p *PolygonApi) fixValues(method string, values url.Values) url.Values {
	tm := time.Now().Unix()
	values.Set("time", strconv.FormatInt(tm, 10))
	values.Set("apiKey", p.confFile.ApiKey)

	builder := bytes.Buffer{}
	rand := "000000"
	builder.WriteString(rand)
	builder.WriteRune('/')
	builder.WriteString(method)
	builder.WriteRune('?')
	builder.WriteString(values.Encode())
	builder.WriteRune('#')
	builder.WriteString(p.confFile.ApiSecret)

	endoded := sha512.Sum512(builder.Bytes())

	values.Set("apiSig", rand+hex.EncodeToString(endoded[:]))
	return values
}
