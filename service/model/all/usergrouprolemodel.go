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
	userGroupRoleFieldNames          = builder.RawFieldNames(&UserGroupRole{})
	userGroupRoleRows                = strings.Join(userGroupRoleFieldNames, ",")
	userGroupRoleRowsExpectAutoSet   = strings.Join(stringx.Remove(userGroupRoleFieldNames, "`create_time`", "`update_time`"), ",")
	userGroupRoleRowsWithPlaceHolder = strings.Join(stringx.Remove(userGroupRoleFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	userGroupRoleListRows            = strings.Join(builder.RawFieldNames(&UserGroupRoleItem{}), ",")

	cacheUserGroupRoleIdPrefix                = "cache:userGroupRole:id:"
	cacheUserGroupRoleUserGroupIdRoleIdPrefix = "cache:userGroupRole:userGroupId:roleId:"
)

type (
	UserGroupRoleModel interface {
		Insert(session sqlx.Session, data *UserGroupRole) (sql.Result, error)
		// Insert(data *UserGroupRole) (sql.Result,error)
		FindOne(id int64) (*UserGroupRole, error)
		FindAll(in *tool.GetsReq) ([]*UserGroupRoleItem, error)
		FindOneByUserGroupIdRoleId(userGroupId int64, roleId int64) (*UserGroupRole, error)
		Update(session sqlx.Session, data *UserGroupRole) error
		// Update(data *UserGroupRole) error
		SoftDelete(session sqlx.Session, data *UserGroupRole) error
		Delete(id int64) error
	}

	defaultUserGroupRoleModel struct {
		sqlc.CachedConn
		table string
	}

	UserGroupRole struct {
		Id          int64     `db:"id"`            // 自增主键
		UserGroupId int64     `db:"user_group_id"` // 用户、组关系表id
		RoleId      int64     `db:"role_id"`       // 角色id
		Status      int64     `db:"status"`        // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime  time.Time `db:"create_time"`   // 创建时间
		UpdateTime  time.Time `db:"update_time"`   // 更新时间
	}
)

func NewUserGroupRoleModel(conn sqlx.SqlConn, c cache.CacheConf) UserGroupRoleModel {
	return &defaultUserGroupRoleModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_group_role`",
	}
}

/*
func (m *defaultUserGroupRoleModel) Insert(data *UserGroupRole) (sql.Result,error) {
	userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, data.Id)
userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, data.UserGroupId, data.RoleId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, userGroupRoleRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserGroupId, data.RoleId, data.Status)
	}, userGroupRoleIdKey, userGroupRoleUserGroupIdRoleIdKey)
	return ret,err
}
*/

func (m *defaultUserGroupRoleModel) Insert(session sqlx.Session, data *UserGroupRole) (sql.Result, error) {
	userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, data.Id)
	userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, data.UserGroupId, data.RoleId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, userGroupRoleRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserGroupId, data.RoleId, data.Status)
		}
		return conn.Exec(query, data.Id, data.UserGroupId, data.RoleId, data.Status)
	}, userGroupRoleIdKey, userGroupRoleUserGroupIdRoleIdKey)
	return ret, err
}

func (m *defaultUserGroupRoleModel) FindOne(id int64) (*UserGroupRole, error) {
	userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, id)
	var resp UserGroupRole
	err := m.QueryRow(&resp, userGroupRoleIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userGroupRoleRows, m.table)
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

func (m *defaultUserGroupRoleModel) FindOneByUserGroupIdRoleId(userGroupId int64, roleId int64) (*UserGroupRole, error) {
	userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, userGroupId, roleId)
	var resp UserGroupRole
	err := m.QueryRowIndex(&resp, userGroupRoleUserGroupIdRoleIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_group_id` = ? and `role_id` = ? limit 1", userGroupRoleRows, m.table)
		if err := conn.QueryRow(&resp, query, userGroupId, roleId); err != nil {
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
func (m *defaultUserGroupRoleModel) Update(data *UserGroupRole) error {
	userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, data.UserGroupId, data.RoleId)
userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, data.Id)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userGroupRoleRowsWithPlaceHolder)
		return conn.Exec(query, data.UserGroupId, data.RoleId, data.Status, data.Id)
	}, userGroupRoleIdKey, userGroupRoleUserGroupIdRoleIdKey)
	return err
}
*/

func (m *defaultUserGroupRoleModel) Update(session sqlx.Session, data *UserGroupRole) error {
	userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, data.UserGroupId, data.RoleId)
	userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userGroupRoleRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserGroupId, data.RoleId, data.Status, data.Id)
		}
		return conn.Exec(query, data.UserGroupId, data.RoleId, data.Status, data.Id)
	}, userGroupRoleIdKey, userGroupRoleUserGroupIdRoleIdKey)
	return err
}

func (m *defaultUserGroupRoleModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userGroupRoleIdKey := fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, id)
	userGroupRoleUserGroupIdRoleIdKey := fmt.Sprintf("%s%v:%v", cacheUserGroupRoleUserGroupIdRoleIdPrefix, data.UserGroupId, data.RoleId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userGroupRoleIdKey, userGroupRoleUserGroupIdRoleIdKey)
	return err
}

func (m *defaultUserGroupRoleModel) SoftDelete(session sqlx.Session, data *UserGroupRole) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultUserGroupRoleModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserGroupRoleIdPrefix, primary)
}

func (m *defaultUserGroupRoleModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userGroupRoleRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultUserGroupRoleModel) FindAll(in *tool.GetsReq) ([]*UserGroupRoleItem, error) {
	resp := make([]*UserGroupRoleItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, userGroupRoleListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
