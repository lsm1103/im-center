package requestsAiDph

//返回该算法的uid
func (r *ReqAiDph) GetAiUidByName(name string) string {
	ais, err := r.GetAisCeche(&GetAisReq{
		Current:  1,
		PageSize: 10,
	})
	if err != nil { return "" }
	if _,ok := ais.AiMap[name]; !ok {return "" }
	return ais.AiMap[name].Uid
}
