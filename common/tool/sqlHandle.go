package tool
//
//import (
//	"fmt"
//	"im-center/api/internal/types"
//	"strings"
//)
//
//type (
//	ModelTool interface {
//		BuildQuery(in *GetsReq, ListRows string, tableName string) string
//	}
//
//	defaultModelTool struct {}
//
//	GetsQueryItem struct {
//		Key        string `json:"key"`                                   // key
//		Val        string `json:"val"`                                   // val int/str('xxx')/list('xx,xx,xx')
//		Handle     string `json:"handle"`                                // 操作方法 =/>/</like
//		NextHandle string `json:"nextHandle,options=or|and,default=and"` // 与下一个条件的操作方法 or/and
//	}
//
//	GetsReq types.GetsReq
//	//GetsReq struct {
//	//	Query    []*GetsQueryItem `json:"query"`                              //查询
//	//	OrderBy  string           `form:"orderBy,default=id"`                 //分组
//	//	Sort     string           `form:"sort,options=desc|asc,default=desc"` //排序 desc/asc
//	//	Current  int64            `form:"current,default=1"`                  //当前页
//	//	PageSize int64            `form:"pageSize,default=10"`                //页面大小
//	//}
//)
//
//func NewModelTool() ModelTool {
//	return &defaultModelTool{}
//}
//
//func (t defaultModelTool) BuildQuery(in *GetsReq, ListRows string, tableName string) string {
//	filter_ := ""
//	tmp_ := ""
//	if in.Query != nil {
//		for index,value := range in.Query {
//			Key := strings.TrimSpace(value.Key)
//			Val := strings.TrimSpace(value.Val)
//			Handle := strings.TrimSpace(value.Handle)
//			NextHandle := strings.TrimSpace(value.NextHandle)
//			switch Handle {
//			case "between":
//				t_ := strings.Split(Val, ",")
//				tmp_ = fmt.Sprintf(" `%s` between '%s' and '%s' ", Key, t_[0], t_[1])
//			case "in":
//				t_ := ""
//				for _, v := range strings.Split(Val, ",") {
//					if len(t_) == 0 {
//						t_ = fmt.Sprintf("'%s'", strings.ToLower(v))
//					} else {
//						t_ = fmt.Sprintf("%s, '%s'", t_, strings.ToLower(v))
//					}
//				}
//				tmp_ = fmt.Sprintf(" `%s` in (%s) ", Key, t_)
//			default:
//				tmp_ = fmt.Sprintf(" `%s` %s '%s' ", Key, Handle, Val)
//			}
//			if index == 0 {
//				filter_ = tmp_
//			} else {
//				filter_ = filter_ + NextHandle + tmp_
//			}
//		}
//		if filter_ != "" {
//			filter_ = fmt.Sprintf("where %s", filter_)
//		}
//	}
//	//Fields_ := ""
//	//for _, v_ := range strings.Split(in.Fields, ",") {
//	//	if len(Fields_) == 0 {
//	//		Fields_ = fmt.Sprintf("`%s`", strings.ToLower(v_))
//	//	} else {
//	//		Fields_ = fmt.Sprintf("%s, `%s`", Fields_, strings.ToLower(v_))
//	//	}
//	//}
//	return fmt.Sprintf(
//		"select %s from %s %s order by %s %v limit %v offset %v",
//		ListRows,
//		tableName,
//		filter_,
//		in.OrderBy,
//		in.Sort,
//		in.PageSize+1,
//		(in.Current - 1) * in.PageSize,
//	)
//}
