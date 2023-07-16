package midware

import (
	"brillinkmicros/internal/conf"
	"brillinkmicros/internal/data"
	"brillinkmicros/pkg"
	"context"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

//const BlAuthHeaderKey = "BL-AUTH-DATA"

func BlAuth(c *conf.Data, dt *data.Data) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				oa := BlPortalOauth2{
					url:            c.BlAuth.Url,
					pathCheckToken: c.BlAuth.PathCheckToken,
				}
				rawToken := tr.RequestHeader().Get("Authorization")
				re := regexp.MustCompile(`Bearer\s+(\S+)`)
				match := re.FindStringSubmatch(rawToken)
				if len(match) != 2 {
					return nil, errors.New(401, "invalid token", rawToken)
				}
				token := match[1]
				var authData []byte
				authData, err = oa.CheckToken(token, c.BlAuth.ClientId, c.BlAuth.ClientSecret)
				if err != nil {
					return nil, err
				}

				dsi, err := getScopesByAuthData(dt, authData)
				if err != nil {
					return nil, errors.Newf(500, "Error Occ At DataScopeServer", err.Error())
				}

				var dataScope []byte
				dataScope, err = json.Marshal(dsi)

				//tr.RequestHeader().Set(BlAuthHeaderKey, string(data))
				tr.RequestHeader().Set(pkg.BlDataScopeHeaderKey, string(dataScope))

				// Do something on entering
				//defer func() {
				//	// Do something on exiting
				//}()
			}
			return handler(ctx, req)
		}
	}
}

type BlPortalOauth2 struct {
	url            string
	pathCheckToken string
}

func (t BlPortalOauth2) CheckToken(token string, clientId string, clientSecret string) ([]byte, error) {
	reqParams := url.Values{}
	reqParams.Set("token", token)
	reqParams.Set("client_id", clientId)
	reqParams.Set("client_secret", clientSecret)
	req, err := http.NewRequest("POST", t.url+t.pathCheckToken, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = reqParams.Encode()
	req.Header.Set("Accept", "*/*")
	cli := &http.Client{Timeout: 10 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	code, err := jsonparser.GetInt(respBody, "code")
	if err != nil {
		return nil, err
	}
	msg, err := jsonparser.GetString(respBody, "msg")
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return nil, errors.New(int(code), "invalid token", msg)
	}

	data, dt, _, err := jsonparser.Get(respBody, "data")
	if err != nil {
		return nil, err
	}
	if dt != jsonparser.ValueType(3) {
		return nil, errors.New(401, "invalid token", msg)
	}
	return data, nil
}

//{"code":401,"data":null,"msg":"访问令牌不存在"}
//{"code":0,"data":{"scopes":null,"exp":1687949113,"user_id":1,"user_type":2,"tenant_id":null,"client_id":"default","access_token":"0e0a56d4dfb2469f9be5a2de469ff59f"},"msg":""}
