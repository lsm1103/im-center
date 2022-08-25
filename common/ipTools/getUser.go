package ipTools

import (
	"net/http"
	"strings"
)
/*获取用户ip
	ip := r.Header.Get("x-forwarded-for")
	ip = r.Header.Get("X-Forwarded-For")
	ip = r.Header.Get("Proxy-Client-IP")
	ip = r.Header.Get("WL-Proxy-Client-IP")
	ip = r.Header.Get("HTTP_CLIENT_IP")
	ip = r.Header.Get("HTTP_X_FORWARDED_FOR")
	ip = r.Header.Get("X-Real-IP")
	ip = r.RemoteAddr
*/
func GetUserIp(r *http.Request) string {
	ip := "0.0.0.0"
	ip_for := []string{"x-forwarded-for", "X-Forwarded-For", "Proxy-Client-IP", "WL-Proxy-Client-IP", "HTTP_CLIENT_IP","HTTP_X_FORWARDED_FOR","X-Real-IP", "RemoteAddr"}
	for _,item := range ip_for {
		if item == "RemoteAddr" {
			tmpIp := strings.Split(r.RemoteAddr,":")[0]
			if len(ip) < 8{
				ip = tmpIp
			}
		} else {
			ip = r.Header.Get(item)
		}
		//fmt.Printf("GetIp-%s: %s\n", item,ip)
		if ip != "" || len(ip) != 0 {
			return ip
		}
	}
	return ip
}
