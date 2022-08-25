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
	thirdAppUseUserFieldNames          = builder.RawFieldNames(&ThirdAppUseUser{})
	thirdAppUseUserRows                = strings.Join(thirdAppUseUserFieldNames, ",")
	thirdAppUseUserRowsExpectAutoSet   = strings.Join(stringx.Remove(thirdAppUseUserFieldNames, "`create_time`", "`update_time`"), ",")
	thirdAppUseUserRowsWithPlaceHolder = strings.Join(stringx.Remove(thirdAppUseUserFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	thirdAppUseUserListRows            = strings.Join(builder.RawFieldNames(&ThirdAppUseUserItem{}), ",")

	cacheThirdAppUseUserIdPrefix                  = "cache:thirdAppUseUser:id:"
	cacheThirdAppUseUserUserIdThirdAppUseIdPrefix = "cache:thirdAppUseUser:userId:thirdAppUseId:"
)

type (
	ThirdAppUseUserModel interface {
		Insert(session sqlx.Session, data *ThirdAppUseUser) (sql.Result, error)
		// Insert(data *ThirdAppUseUser) (sql.Result,error)
		FindOne(id int64) (*ThirdAppUseUser, error)
		FindAll(in *tool.GetsReq) ([]*ThirdAppUseUserItem, error)
		FindOneByUserIdThirdAppUseId(userId int64, thirdAppUseId int64) (*ThirdAppUseUser, error)
		Update(session sqlx.Session, data *ThirdAppUseUser) error
		// Update(data *ThirdAppUseUser) error
		SoftDelete(session sqlx.Session, data *ThirdAppUseUser) error
		Delete(id int64) error
	}

	defaultThirdAppUseUserModel struct {
		sqlc.CachedConn
		table string
	}

	ThirdAppUseUser struct {
		Id            int64     `db:"id"`              // 自增主键，open_id
		ThirdAppUseId int64     `db:"thirdApp_use_id"` // 第三方授权应用唯一id
		UserId        int64     `db:"user_id"`         // 用户id
		AuthScope     string    `db:"auth_scope"`      // 授权范围
		Extra         string    `db:"extra"`           // 附加属性
		Status        int64     `db:"status"`          // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime    time.Time `db:"create_time"`     // 创建时间
		UpdateTime    time.Time `db:"update_time"`     // 更新时间
	}
)

func NewThirdAppUseUserModel(conn sqlx.SqlConn, c cache.CacheConf) ThirdAppUseUserModel {
	return &defaultThirdAppUseUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`thirdApp_use_user`",
	}
}

/*
func (m *defaultThirdAppUseUserModel) Insert(data *ThirdAppUseUser) (sql.Result,error) {
	thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, data.Id)
thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, data.UserId, data.ThirdAppUseId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, thirdAppUseUserRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status)
	}, thirdAppUseUserIdKey, thirdAppUseUserUserIdThirdAppUseIdKey)
	return ret,err
}
*/

func (m *defaultThirdAppUseUserModel) Insert(session sqlx.Session, data *ThirdAppUseUser) (sql.Result, error) {
	thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, data.Id)
	thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, data.UserId, data.ThirdAppUseId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, thirdAppUseUserRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status)
		}
		return conn.Exec(query, data.Id, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status)
	}, thirdAppUseUserIdKey, thirdAppUseUserUserIdThirdAppUseIdKey)
	return ret, err
}

func (m *defaultThirdAppUseUserModel) FindOne(id int64) (*ThirdAppUseUser, error) {
	thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, id)
	var resp ThirdAppUseUser
	err := m.QueryRow(&resp, thirdAppUseUserIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdAppUseUserRows, m.table)
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

func (m *defaultThirdAppUseUserModel) FindOneByUserIdThirdAppUseId(userId int64, thirdAppUseId int64) (*ThirdAppUseUser, error) {
	thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, userId, thirdAppUseId)
	var resp ThirdAppUseUser
	err := m.QueryRowIndex(&resp, thirdAppUseUserUserIdThirdAppUseIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `thirdApp_use_id` = ? limit 1", thirdAppUseUserRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, thirdAppUseId); err != nil {
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
func (m *defaultThirdAppUseUserModel) Update(data *ThirdAppUseUser) error {
	thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, data.UserId, data.ThirdAppUseId)
thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, data.Id)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdAppUseUserRowsWithPlaceHolder)
		return conn.Exec(query, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status, data.Id)
	}, thirdAppUseUserIdKey, thirdAppUseUserUserIdThirdAppUseIdKey)
	return err
}
*/

func (m *defaultThirdAppUseUserModel) Update(session sqlx.Session, data *ThirdAppUseUser) error {
	thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, data.UserId, data.ThirdAppUseId)
	thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, thirdAppUseUserRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status, data.Id)
		}
		return conn.Exec(query, data.ThirdAppUseId, data.UserId, data.AuthScope, data.Extra, data.Status, data.Id)
	}, thirdAppUseUserIdKey, thirdAppUseUserUserIdThirdAppUseIdKey)
	return err
}

func (m *defaultThirdAppUseUserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	thirdAppUseUserIdKey := fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, id)
	thirdAppUseUserUserIdThirdAppUseIdKey := fmt.Sprintf("%s%v:%v", cacheThirdAppUseUserUserIdThirdAppUseIdPrefix, data.UserId, data.ThirdAppUseId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, thirdAppUseUserIdKey, thirdAppUseUserUserIdThirdAppUseIdKey)
	return err
}

func (m *defaultThirdAppUseUserModel) SoftDelete(session sqlx.Session, data *ThirdAppUseUser) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultThirdAppUseUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheThirdAppUseUserIdPrefix, primary)
}

func (m *defaultThirdAppUseUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", thirdAppUseUserRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultThirdAppUseUserModel) FindAll(in *tool.GetsReq) ([]*ThirdAppUseUserItem, error) {
	resp := make([]*ThirdAppUseUserItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, thirdAppUseUserListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
