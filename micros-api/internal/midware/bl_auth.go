package midware

import (
	"context"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/redis/go-redis/v9"
	"io"
	"micros-api/internal/conf"
	"micros-api/internal/data"
	"micros-api/pkg"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

//const BlAuthHeaderKey = "BL-AUTH-DATA"

var oAuth *OAuth

func BlAuth(c *conf.Data, dt *data.Data) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {

				rawToken := tr.RequestHeader().Get("Authorization")
				re := regexp.MustCompile(`Bearer\s+(\S+)`)
				match := re.FindStringSubmatch(rawToken)
				if len(match) != 2 {
					return nil, errors.New(401, "invalid token", rawToken)
				}

				if oAuth == nil {
					oAuth = NewOAuth(
						c.BlAuth.ClientId,
						c.BlAuth.ClientSecret,
						c.BlAuth.Url,
						c.BlAuth.PathCheckToken,
					)
				}

				token := match[1]
				var authData []byte
				authData, err = oAuth.CheckToken(token)
				if err != nil {
					return nil, err
				}
				userId, err := oAuth.parseUserId(authData)
				if err != nil {
					return nil, err
				}

				// cache
				dsi, err := oAuth.getCachedScopesByUserId(userId, dt)
				if err != nil {
					return nil, errors.Newf(500, "Error Occ At DataScopeServer:oAuth.getCachedScopesByUserId: %v", err.Error())
				}
				if dsi == nil {
					dsi, err = getScopesByUserId(dt, userId)
					if err != nil {
						return nil, errors.Newf(500, "Error Occ At DataScopeServer:getScopesByUserId: %v", err.Error())
					}
					err = oAuth.setCachedScopesByUserId(userId, dsi, dt)
					if err != nil {
						return nil, errors.Newf(500, "Error Occ At DataScopeServer:oAuth.setCachedScopesByUserId: %v", err.Error())
					}
				}
				// cache

				//dsi, err := getScopesByUserId(dt, userId)
				if err != nil {
					return nil, errors.Newf(500, "Error Occ At DataScopeServer: %v", err.Error())
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

type OAuth struct {
	clientId     string
	clientSecret string
	apiUrl       string
}

func NewOAuth(clientId, clientSecret, url, tokenCheckPath string) *OAuth {
	return &OAuth{
		clientId:     clientId,
		clientSecret: clientSecret,
		apiUrl:       url + tokenCheckPath,
	}
}

func (o *OAuth) dataScopeKey(userId int64) string {
	return "data_scopes:" + strconv.FormatInt(userId, 10)
}

func (o *OAuth) getCachedScopesByUserId(userId int64, dt *data.Data) (*pkg.DataScopeInfo, error) {
	key := o.dataScopeKey(userId)
	b, err := dt.Rdb.Client.Get(context.TODO(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			// Handle the timeout error
			return nil, nil
		}
		return nil, err
	}
	var dsi pkg.DataScopeInfo
	err = json.Unmarshal(b, &dsi)
	if err != nil {
		return nil, err
	}
	return &dsi, nil
}

func (o *OAuth) setCachedScopesByUserId(userId int64, dsi *pkg.DataScopeInfo, dt *data.Data) error {
	key := o.dataScopeKey(userId)
	b, err := json.Marshal(dsi)
	if err != nil {
		return err
	}
	err = dt.Rdb.Client.Set(context.TODO(), key, b, time.Second*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func (o *OAuth) parseUserId(authRes []byte) (int64, error) {
	var userId int64
	//userId, err := jsonparser.GetInt(authData, "user_id")
	userIdBytes, dType, _, err := jsonparser.Get(authRes, "user_id")
	if err != nil {
		return 0, err
	}
	switch dType {
	case jsonparser.String:
		userId, err = strconv.ParseInt(string(userIdBytes), 10, 64)
		if err != nil {
			return 0, err
		}
	case jsonparser.Number:
		userId, err = jsonparser.ParseInt(userIdBytes)
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.New(500, "At parseScopes dtype of user_id upexpect", string(userIdBytes))
	}

	return userId, nil
}

func (o *OAuth) CheckToken(token string) ([]byte, error) {
	reqParams := url.Values{}
	reqParams.Set("token", token)
	reqParams.Set("client_id", o.clientId)
	reqParams.Set("client_secret", o.clientSecret)
	req, err := http.NewRequest("POST", o.apiUrl, nil)
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

	authData, dt, _, err := jsonparser.Get(respBody, "data")
	if err != nil {
		return nil, err
	}
	if dt != jsonparser.ValueType(3) {
		return nil, errors.New(401, "invalid token", msg)
	}
	return authData, nil
}

//{"code":401,"data":null,"msg":"访问令牌不存在"}
//{"code":0,"data":{"scopes":null,"exp":1687949113,"user_id":1,"user_type":2,"tenant_id":null,"client_id":"default","access_token":"0e0a56d4dfb2469f9be5a2de469ff59f"},"msg":""}
