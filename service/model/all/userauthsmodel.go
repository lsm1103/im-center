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
	userAuthsFieldNames          = builder.RawFieldNames(&UserAuths{})
	userAuthsRows                = strings.Join(userAuthsFieldNames, ",")
	userAuthsRowsExpectAutoSet   = strings.Join(stringx.Remove(userAuthsFieldNames, "`create_time`", "`update_time`"), ",")
	userAuthsRowsWithPlaceHolder = strings.Join(stringx.Remove(userAuthsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	userAuthsListRows            = strings.Join(builder.RawFieldNames(&UserAuthsItem{}), ",")

	cacheUserAuthsIdPrefix         = "cache:userAuths:id:"
	cacheUserAuthsIdentifierPrefix = "cache:userAuths:identifier:"
)

type (
	UserAuthsModel interface {
		Insert(session sqlx.Session, data *UserAuths) (sql.Result, error)
		// Insert(data *UserAuths) (sql.Result,error)
		FindOne(id int64) (*UserAuths, error)
		FindAll(in *tool.GetsReq) ([]*UserAuthsItem, error)
		FindOneByIdentifier(identifier string) (*UserAuths, error)
		Update(session sqlx.Session, data *UserAuths) error
		// Update(data *UserAuths) error
		SoftDelete(session sqlx.Session, data *UserAuths) error
		Delete(id int64) error
	}

	defaultUserAuthsModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuths struct {
		Id           int64     `db:"id"`            // 自增主键
		UserId       int64     `db:"user_id"`       // 用户id
		IdentityType string    `db:"identity_type"` // 身份认证方式，手机号/邮箱/身份证/第三方
		Identifier   string    `db:"identifier"`    // 唯一标识符，手机号/邮箱/身份证/第三方open_id
		Status       int64     `db:"status"`        // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime   time.Time `db:"create_time"`   // 创建时间
		UpdateTime   time.Time `db:"update_time"`   // 更新时间
	}
)

func NewUserAuthsModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthsModel {
	return &defaultUserAuthsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auths`",
	}
}

/*
func (m *defaultUserAuthsModel) Insert(data *UserAuths) (sql.Result,error) {
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, data.Id)
userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, data.Identifier)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userAuthsRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserId, data.IdentityType, data.Identifier, data.Status)
	}, userAuthsIdKey, userAuthsIdentifierKey)
	return ret,err
}
*/

func (m *defaultUserAuthsModel) Insert(session sqlx.Session, data *UserAuths) (sql.Result, error) {
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, data.Id)
	userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, data.Identifier)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userAuthsRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserId, data.IdentityType, data.Identifier, data.Status)
		}
		return conn.Exec(query, data.Id, data.UserId, data.IdentityType, data.Identifier, data.Status)
	}, userAuthsIdKey, userAuthsIdentifierKey)
	return ret, err
}

func (m *defaultUserAuthsModel) FindOne(id int64) (*UserAuths, error) {
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, id)
	var resp UserAuths
	err := m.QueryRow(&resp, userAuthsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthsRows, m.table)
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

func (m *defaultUserAuthsModel) FindOneByIdentifier(identifier string) (*UserAuths, error) {
	userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, identifier)
	var resp UserAuths
	err := m.QueryRowIndex(&resp, userAuthsIdentifierKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `identifier` = ? limit 1", userAuthsRows, m.table)
		if err := conn.QueryRow(&resp, query, identifier); err != nil {
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
func (m *defaultUserAuthsModel) Update(data *UserAuths) error {
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, data.Id)
userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, data.Identifier)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthsRowsWithPlaceHolder)
		return conn.Exec(query, data.UserId, data.IdentityType, data.Identifier, data.Status, data.Id)
	}, userAuthsIdKey, userAuthsIdentifierKey)
	return err
}
*/

func (m *defaultUserAuthsModel) Update(session sqlx.Session, data *UserAuths) error {
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, data.Id)
	userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, data.Identifier)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthsRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserId, data.IdentityType, data.Identifier, data.Status, data.Id)
		}
		return conn.Exec(query, data.UserId, data.IdentityType, data.Identifier, data.Status, data.Id)
	}, userAuthsIdKey, userAuthsIdentifierKey)
	return err
}

func (m *defaultUserAuthsModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userAuthsIdKey := fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, id)
	userAuthsIdentifierKey := fmt.Sprintf("%s%v", cacheUserAuthsIdentifierPrefix, data.Identifier)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userAuthsIdKey, userAuthsIdentifierKey)
	return err
}

func (m *defaultUserAuthsModel) SoftDelete(session sqlx.Session, data *UserAuths) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultUserAuthsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserAuthsIdPrefix, primary)
}

func (m *defaultUserAuthsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthsRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultUserAuthsModel) FindAll(in *tool.GetsReq) ([]*UserAuthsItem, error) {
	resp := make([]*UserAuthsItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, userAuthsListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
