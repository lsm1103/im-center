package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/token"
	"im-center/common/ctxdata"
)

//生成Jwt令牌
func GenerateJwtToken(secretKey string, iat int64, seconds int64, userData map[string]string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["UserData"] = userData
	access_token := jwt.New(jwt.SigningMethodHS256)
	access_token.Claims = claims
	return access_token.SignedString([]byte(secretKey))
}

//解析验证Jwt令牌
func ParseVerifyJwtToken(r *http.Request, accessSecret string) (map[string]string, error) {
	parser := token.NewTokenParser()
	tok, err := parser.ParseToken(r, accessSecret, "")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ParseVerifyJwtToken Err:%v", err))
	}
	if !tok.Valid {
		return nil, errors.New("ParseVerifyJwtToken Err:令牌已经失效")
	}
	claims, ok := tok.Claims.(jwt.MapClaims) // 解析token中内容
	if !ok {
		return nil, errors.New(fmt.Sprintf("tok.Claims is not ok, tok.Claims:%+v , claims:%+v, ok:%v\n", tok.Claims, claims, ok))
	}
	UserData := map[string]string{}
	//strRet, err := json.Marshal(claims[ctxdata.CtxKeyUserData])
	//err = json.Unmarshal(strRet, UserData)
	if item, ok := claims[ctxdata.CtxKeyUserData].(map[string]interface{}); ok {
		for k, v := range item {
			if val, ok := v.(string); ok {
				UserData[k] = val
			}
		}
	}
	return UserData,err
}

func ParseJwtToken(tokenStr string, accessSecret string) (map[string]string, error){
	tok, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims) // 解析token中内容
	if !ok {
		return nil, errors.New(fmt.Sprintf("tok.Claims is not ok, tok.Claims:%+v , claims:%+v, ok:%v\n", tok.Claims, claims, ok))
	}
	UserData := map[string]string{}
	//strRet, err := json.Marshal(claims[ctxdata.CtxKeyUserData])
	//err = json.Unmarshal(strRet, UserData)
	if item, ok := claims[ctxdata.CtxKeyUserData].(map[string]interface{}); ok {
		for k, v := range item {
			if val, ok := v.(string); ok {
				UserData[k] = val
			}
		}
	}
	return UserData,err
}