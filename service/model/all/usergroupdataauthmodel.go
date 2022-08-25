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
	userGroupDataAuthFieldNames          = builder.RawFieldNames(&UserGroupDataAuth{})
	userGroupDataAuthRows                = strings.Join(userGroupDataAuthFieldNames, ",")
	userGroupDataAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(userGroupDataAuthFieldNames, "`create_time`", "`update_time`"), ",")
	userGroupDataAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(userGroupDataAuthFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	userGroupDataAuthListRows            = strings.Join(builder.RawFieldNames(&UserGroupDataAuthItem{}), ",")

	cacheUserGroupDataAuthIdPrefix                            = "cache:userGroupDataAuth:id:"
	cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix = "cache:userGroupDataAuth:userGroupId:resourceIdentifier:"
)

type (
	UserGroupDataAuthModel interface {
		Insert(session sqlx.Session, data *UserGroupDataAuth) (sql.Result, error)
		// Insert(data *UserGroupDataAuth) (sql.Result,error)
		FindOne(id int64) (*UserGroupDataAuth, error)
		FindAll(in *tool.GetsReq) ([]*UserGroupDataAuthItem, error)
		FindOneByUserGroupIdResourceIdentifier(userGroupId int64, resourceIdentifier string) (*UserGroupDataAuth, error)
		Update(session sqlx.Session, data *UserGroupDataAuth) error
		// Update(data *UserGroupDataAuth) error
		SoftDelete(session sqlx.Session, data *UserGroupDataAuth) error
		Delete(id int64) error
	}

	defaultUserGroupDataAuthModel struct {
		sqlc.CachedConn
		table string
	}

	UserGroupDataAuth struct {
		Id                 int64     `db:"id"`                  // 自增主键
		UserGroupId        int64     `db:"user_group_id"`       // 用户、组关系表id
		ResourceIdentifier string    `db:"resource_identifier"` // 资源标识
		Auth               string    `db:"auth"`                // 权限（all、只看自己、看协同、看同级、看同组别、看同部门、看同公司(机构)、看同产品）
		Status             int64     `db:"status"`              // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime         time.Time `db:"create_time"`         // 创建时间
		UpdateTime         time.Time `db:"update_time"`         // 更新时间
	}
)

func NewUserGroupDataAuthModel(conn sqlx.SqlConn, c cache.CacheConf) UserGroupDataAuthModel {
	return &defaultUserGroupDataAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_group_data_auth`",
	}
}

/*
func (m *defaultUserGroupDataAuthModel) Insert(data *UserGroupDataAuth) (sql.Result,error) {
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, data.Id)
userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, data.UserGroupId, data.ResourceIdentifier)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userGroupDataAuthRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status)
	}, userGroupDataAuthIdKey, userGroupDataAuthUserGroupIdResourceIdentifierKey)
	return ret,err
}
*/

func (m *defaultUserGroupDataAuthModel) Insert(session sqlx.Session, data *UserGroupDataAuth) (sql.Result, error) {
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, data.Id)
	userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, data.UserGroupId, data.ResourceIdentifier)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userGroupDataAuthRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status)
		}
		return conn.Exec(query, data.Id, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status)
	}, userGroupDataAuthIdKey, userGroupDataAuthUserGroupIdResourceIdentifierKey)
	return ret, err
}

func (m *defaultUserGroupDataAuthModel) FindOne(id int64) (*UserGroupDataAuth, error) {
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, id)
	var resp UserGroupDataAuth
	err := m.QueryRow(&resp, userGroupDataAuthIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userGroupDataAuthRows, m.table)
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

func (m *defaultUserGroupDataAuthModel) FindOneByUserGroupIdResourceIdentifier(userGroupId int64, resourceIdentifier string) (*UserGroupDataAuth, error) {
	userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, userGroupId, resourceIdentifier)
	var resp UserGroupDataAuth
	err := m.QueryRowIndex(&resp, userGroupDataAuthUserGroupIdResourceIdentifierKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_group_id` = ? and `resource_identifier` = ? limit 1", userGroupDataAuthRows, m.table)
		if err := conn.QueryRow(&resp, query, userGroupId, resourceIdentifier); err != nil {
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
func (m *defaultUserGroupDataAuthModel) Update(data *UserGroupDataAuth) error {
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, data.Id)
userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, data.UserGroupId, data.ResourceIdentifier)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userGroupDataAuthRowsWithPlaceHolder)
		return conn.Exec(query, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status, data.Id)
	}, userGroupDataAuthIdKey, userGroupDataAuthUserGroupIdResourceIdentifierKey)
	return err
}
*/

func (m *defaultUserGroupDataAuthModel) Update(session sqlx.Session, data *UserGroupDataAuth) error {
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, data.Id)
	userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, data.UserGroupId, data.ResourceIdentifier)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userGroupDataAuthRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status, data.Id)
		}
		return conn.Exec(query, data.UserGroupId, data.ResourceIdentifier, data.Auth, data.Status, data.Id)
	}, userGroupDataAuthIdKey, userGroupDataAuthUserGroupIdResourceIdentifierKey)
	return err
}

func (m *defaultUserGroupDataAuthModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userGroupDataAuthIdKey := fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, id)
	userGroupDataAuthUserGroupIdResourceIdentifierKey := fmt.Sprintf("%s%v:%v", cacheUserGroupDataAuthUserGroupIdResourceIdentifierPrefix, data.UserGroupId, data.ResourceIdentifier)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userGroupDataAuthIdKey, userGroupDataAuthUserGroupIdResourceIdentifierKey)
	return err
}

func (m *defaultUserGroupDataAuthModel) SoftDelete(session sqlx.Session, data *UserGroupDataAuth) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultUserGroupDataAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserGroupDataAuthIdPrefix, primary)
}

func (m *defaultUserGroupDataAuthModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userGroupDataAuthRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultUserGroupDataAuthModel) FindAll(in *tool.GetsReq) ([]*UserGroupDataAuthItem, error) {
	resp := make([]*UserGroupDataAuthItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, userGroupDataAuthListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
