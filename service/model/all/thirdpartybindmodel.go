package all

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"

	"im-center/common/globalkey"
	"im-center/common/tool"
)

var (
	thirdpartyBindFieldNames          = builder.RawFieldNames(&ThirdpartyBind{})
	thirdpartyBindRows                = strings.Join(thirdpartyBindFieldNames, ",")
	thirdpartyBindRowsExpectAutoSet   = strings.Join(stringx.Remove(thirdpartyBindFieldNames, "`create_time`", "`update_time`"), ",")
	thirdpartyBindRowsWithPlaceHolder = strings.Join(stringx.Remove(thirdpartyBindFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	thirdpartyBindListRows            = strings.Join(builder.RawFieldNames(&ThirdpartyBindItem{}), ",")

	cacheThirdpartyBindIdPrefix     = "cache:thirdpartyBind:id:"
	cacheThirdpartyBindOpenIdPrefix = "cache:thirdpartyBind:openId:"
)

type (
	ThirdpartyBindModel interface {
		Insert(session sqlx.Session, data *ThirdpartyBind) (sql.Result, error)
		// Insert(data *ThirdpartyBind) (sql.Result,error)
		FindOne(id int64) (*ThirdpartyBind, error)
		FindAll(in *tool.GetsReq) ([]*ThirdpartyBindItem, error)
		FindOneByOpenId(openId string) (*ThirdpartyBind, error)
		Update(session sqlx.Session, data *ThirdpartyBind) error
		// Update(data *ThirdpartyBind) error
		SoftDelete(session sqlx.Session, data *ThirdpartyBind) error
		Delete(id int64) error
	}

	defaultThirdpartyBindModel struct {
		sqlc.CachedConn
		table string
	}

	ThirdpartyBind struct {
		Id         int64     `db:"id"`          // 自增主键
		UserId     int64     `db:"user_id"`     // 用户id
		OpenId     string    `db:"open_id"`     // 第三方唯一id
		Source     string    `db:"source"`      // 用户来源, 微信/qq/抖音/钉钉/...
		AuthScope  string    `db:"auth_scope"`  // 授权范围
		Account    string    `db:"account"`     // 账号
		Nickname   string    `db:"nickname"`    // 用户昵称
		Gender     string    `db:"gender"`      // 性别
		Avatar     string    `db:"avatar"`      // 用户头像
		Blog       string    `db:"blog"`        // 用户网址
		Company    string    `db:"company"`     // 所在公司
		Location   string    `db:"location"`    // 位置
		Email      string    `db:"email"`       // 邮箱
		Remark     string    `db:"remark"`      // 用户备注（各平台中的用户个人介绍）
		Extra      string    `db:"extra"`       // 附加属性
		Status     int64     `db:"status"`      // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
	}
)

func NewThirdpartyBindModel(conn sqlx.SqlConn, c cache.CacheConf) ThirdpartyBindModel {
	return &defaultThirdpartyBindModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`thirdparty_bind`",
	}
}

/*
func (m *defaultThirdpartyBindModel) Insert(data *ThirdpartyBind) (sql.Result,error) {
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, data.Id)
thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, data.OpenId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, thirdpartyBindRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status)
	}, thirdpartyBindIdKey, thirdpartyBindOpenIdKey)
	return ret,err
}
*/

func (m *defaultThirdpartyBindModel) Insert(session sqlx.Session, data *ThirdpartyBind) (sql.Result, error) {
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, data.Id)
	thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, data.OpenId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, thirdpartyBindRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status)
		}
		return conn.Exec(query, data.Id, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status)
	}, thirdpartyBindIdKey, thirdpartyBindOpenIdKey)
	return ret, err
}

func (m *defaultThirdpartyBindModel) FindOne(id int64) (*ThirdpartyBind, error) {
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, id)
	var resp ThirdpartyBind
	err := m.QueryRow(&resp, thirdpartyBindIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdpartyBindRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultThirdpartyBindModel) FindOneByOpenId(openId string) (*ThirdpartyBind, error) {
	thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, openId)
	var resp ThirdpartyBind
	err := m.QueryRowIndex(&resp, thirdpartyBindOpenIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `open_id` = ? limit 1", thirdpartyBindRows, m.table)
		if err := conn.QueryRow(&resp, query, openId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

/*
func (m *defaultThirdpartyBindModel) Update(data *ThirdpartyBind) error {
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, data.Id)
thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, data.OpenId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdpartyBindRowsWithPlaceHolder)
		return conn.Exec(query, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status, data.Id)
	}, thirdpartyBindIdKey, thirdpartyBindOpenIdKey)
	return err
}
*/

func (m *defaultThirdpartyBindModel) Update(session sqlx.Session, data *ThirdpartyBind) error {
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, data.Id)
	thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, data.OpenId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdpartyBindRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status, data.Id)
		}
		return conn.Exec(query, data.UserId, data.OpenId, data.Source, data.AuthScope, data.Account, data.Nickname, data.Gender, data.Avatar, data.Blog, data.Company, data.Location, data.Email, data.Remark, data.Extra, data.Status, data.Id)
	}, thirdpartyBindIdKey, thirdpartyBindOpenIdKey)
	return err
}

func (m *defaultThirdpartyBindModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	thirdpartyBindIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, id)
	thirdpartyBindOpenIdKey := fmt.Sprintf("%s%v", cacheThirdpartyBindOpenIdPrefix, data.OpenId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, thirdpartyBindOpenIdKey, thirdpartyBindIdKey)
	return err
}

func (m *defaultThirdpartyBindModel) SoftDelete(session sqlx.Session, data *ThirdpartyBind) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultThirdpartyBindModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheThirdpartyBindIdPrefix, primary)
}

func (m *defaultThirdpartyBindModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdpartyBindRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultThirdpartyBindModel) FindAll(in *tool.GetsReq) ([]*ThirdpartyBindItem, error) {
	resp := make([]*ThirdpartyBindItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, thirdpartyBindListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
