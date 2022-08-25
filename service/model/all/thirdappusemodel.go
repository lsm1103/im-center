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
	thirdAppUseFieldNames          = builder.RawFieldNames(&ThirdAppUse{})
	thirdAppUseRows                = strings.Join(thirdAppUseFieldNames, ",")
	thirdAppUseRowsExpectAutoSet   = strings.Join(stringx.Remove(thirdAppUseFieldNames, "`create_time`", "`update_time`"), ",")
	thirdAppUseRowsWithPlaceHolder = strings.Join(stringx.Remove(thirdAppUseFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	thirdAppUseListRows            = strings.Join(builder.RawFieldNames(&ThirdAppUseItem{}), ",")

	cacheThirdAppUseIdPrefix    = "cache:thirdAppUse:id:"
	cacheThirdAppUseAppIdPrefix = "cache:thirdAppUse:appId:"
)

type (
	ThirdAppUseModel interface {
		Insert(session sqlx.Session, data *ThirdAppUse) (sql.Result, error)
		// Insert(data *ThirdAppUse) (sql.Result,error)
		FindOne(id int64) (*ThirdAppUse, error)
		FindAll(in *tool.GetsReq) ([]*ThirdAppUseItem, error)
		FindOneByAppId(appId string) (*ThirdAppUse, error)
		Update(session sqlx.Session, data *ThirdAppUse) error
		// Update(data *ThirdAppUse) error
		SoftDelete(session sqlx.Session, data *ThirdAppUse) error
		Delete(id int64) error
	}

	defaultThirdAppUseModel struct {
		sqlc.CachedConn
		table string
	}

	ThirdAppUse struct {
		Id          int64     `db:"id"`           // 自增主键
		AppId       string    `db:"app_id"`       // app唯一id
		AppSecret   string    `db:"app_secret"`   // app密码
		AuthScope   string    `db:"auth_scope"`   // 授权范围
		Name        string    `db:"name"`         // app名称
		CallbackUrl string    `db:"callback_url"` // 回调url
		Ico         string    `db:"ico"`          // 图标
		Email       string    `db:"email"`        // 邮箱
		Phone       string    `db:"phone"`        // 电话
		Remark      string    `db:"remark"`       // 备注
		Extra       string    `db:"extra"`        // 附加属性
		Status      int64     `db:"status"`       // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 更新时间
	}
)

func NewThirdAppUseModel(conn sqlx.SqlConn, c cache.CacheConf) ThirdAppUseModel {
	return &defaultThirdAppUseModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`thirdApp_use`",
	}
}

/*
func (m *defaultThirdAppUseModel) Insert(data *ThirdAppUse) (sql.Result,error) {
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, data.Id)
thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, data.AppId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, thirdAppUseRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status)
	}, thirdAppUseAppIdKey, thirdAppUseIdKey)
	return ret,err
}
*/

func (m *defaultThirdAppUseModel) Insert(session sqlx.Session, data *ThirdAppUse) (sql.Result, error) {
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, data.Id)
	thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, data.AppId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, thirdAppUseRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status)
		}
		return conn.Exec(query, data.Id, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status)
	}, thirdAppUseAppIdKey, thirdAppUseIdKey)
	return ret, err
}

func (m *defaultThirdAppUseModel) FindOne(id int64) (*ThirdAppUse, error) {
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, id)
	var resp ThirdAppUse
	err := m.QueryRow(&resp, thirdAppUseIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdAppUseRows, m.table)
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

func (m *defaultThirdAppUseModel) FindOneByAppId(appId string) (*ThirdAppUse, error) {
	thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, appId)
	var resp ThirdAppUse
	err := m.QueryRowIndex(&resp, thirdAppUseAppIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `app_id` = ? limit 1", thirdAppUseRows, m.table)
		if err := conn.QueryRow(&resp, query, appId); err != nil {
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
func (m *defaultThirdAppUseModel) Update(data *ThirdAppUse) error {
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, data.Id)
thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, data.AppId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdAppUseRowsWithPlaceHolder)
		return conn.Exec(query, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status, data.Id)
	}, thirdAppUseIdKey, thirdAppUseAppIdKey)
	return err
}
*/

func (m *defaultThirdAppUseModel) Update(session sqlx.Session, data *ThirdAppUse) error {
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, data.Id)
	thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, data.AppId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdAppUseRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status, data.Id)
		}
		return conn.Exec(query, data.AppId, data.AppSecret, data.AuthScope, data.Name, data.CallbackUrl, data.Ico, data.Email, data.Phone, data.Remark, data.Extra, data.Status, data.Id)
	}, thirdAppUseIdKey, thirdAppUseAppIdKey)
	return err
}

func (m *defaultThirdAppUseModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	thirdAppUseIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, id)
	thirdAppUseAppIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseAppIdPrefix, data.AppId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, thirdAppUseIdKey, thirdAppUseAppIdKey)
	return err
}

func (m *defaultThirdAppUseModel) SoftDelete(session sqlx.Session, data *ThirdAppUse) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultThirdAppUseModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheThirdAppUseIdPrefix, primary)
}

func (m *defaultThirdAppUseModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdAppUseRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultThirdAppUseModel) FindAll(in *tool.GetsReq) ([]*ThirdAppUseItem, error) {
	resp := make([]*ThirdAppUseItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, thirdAppUseListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
