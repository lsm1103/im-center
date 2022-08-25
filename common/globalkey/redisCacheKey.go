package globalkey

import "fmt"

/**
redis key except "model cache key"  in here,
but "model cache key" in model
*/

// 用户登陆的token  			userId:deviceId
const CacheUserTokenKey = "token:%+v:%+v:user_token"
// 用户的refresh_token  		userId:deviceId
const CacheRefreshTokenKey = "token:%+v:%+v:refresh_token"

// 用户的OAuth2 code 		userId:appUseId
const CacheOAuth2CodeKey = "oauth2:%+v:%+v:code"
// 用户的OAuth2 token 		userId:appUseId
const CacheOAuth2TokenKey = "oauth2:%+v:%+v:token"
// 用户的OAuth2 refresh_token 	   userId:appUseId
const CacheOAuth2RefreshTokenKey = "oauth2:%+v:%+v:refresh_token"


// 用户的信息by token 		token
const CacheUserInfoByTokenKey = "userInfo:%+v"



// 连接信息的缓存	userId:deviceId
const CacheConnectIdKey = "imConn:userId:deviceId:%s:%s"
// 批量获取所有连接信息key
const CacheConnectIdMATCHKey = "imConn:userId:deviceId:*"
// 批量获取用户的所有连接信息key
const CacheConnectIdUserIdMATCHKey = "imConn:userId:deviceId:%s:*"

// 节点信息的缓存	userId:deviceId
const CacheNodeIdKey = "imConn:nodeId:%s"
const CacheNodeIdMATCHKey = "imConn:nodeId:*"

const ConnectIdFmt = "%s:%s"


func BuildKey(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}