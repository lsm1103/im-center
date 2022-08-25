# 连接层的标准化连接、连接管理、分布式处理



```go

//server.AddRoutes([]rest.Route{
//	{
//		Method:  http.MethodPost,
//		Path:    "/sendOneMsg",
//		Handler: func(w http.ResponseWriter, r *http.Request) {
//			var msg types.SendOneMsgReq
//			err := json.NewDecoder(r.Body).Decode(&msg)
//			if err != nil {
//				fmt.Println("解析参数失败", err)
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//			//fmt.Println("发送消息", msg)
//			err = ctx.Cs.SendOneMsg(&msg)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//			w.Write([]byte([]byte(`{"code":200,"msg":"ok"}`)))
//		},
//	},
//	{
//		Method:  http.MethodGet,
//		Path:    "/getSysInfo",
//		Handler: func(w http.ResponseWriter, r *http.Request) {
//			info, err := ctx.Cs.GetServiceInfo(nil)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
//				return
//			}
//			resp := map[string]interface{}{
//				"code":200,
//				"msg":"ok",
//				"data":info,
//			}
//			d, err := json.Marshal(resp)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusNotFound)
//				return
//			}
//			w.Write(d)
//		},
//	},
//})
```