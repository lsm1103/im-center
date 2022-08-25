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
	roleFieldNames          = builder.RawFieldNames(&Role{})
	roleRows                = strings.Join(roleFieldNames, ",")
	roleRowsExpectAutoSet   = strings.Join(stringx.Remove(roleFieldNames, "`create_time`", "`update_time`"), ",")
	roleRowsWithPlaceHolder = strings.Join(stringx.Remove(roleFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	roleListRows            = strings.Join(builder.RawFieldNames(&RoleItem{}), ",")

	cacheRoleIdPrefix   = "cache:role:id:"
	cacheRoleNamePrefix = "cache:role:name:"
)

type (
	RoleModel interface {
		Insert(session sqlx.Session, data *Role) (sql.Result, error)
		// Insert(data *Role) (sql.Result,error)
		FindOne(id int64) (*Role, error)
		FindAll(in *tool.GetsReq) ([]*RoleItem, error)
		FindOneByName(name string) (*Role, error)
		Update(session sqlx.Session, data *Role) error
		// Update(data *Role) error
		SoftDelete(session sqlx.Session, data *Role) error
		Delete(id int64) error
	}

	defaultRoleModel struct {
		sqlc.CachedConn
		table string
	}

	Role struct {
		Id             int64     `db:"id"`              // 自增主键
		Name           string    `db:"name"`            // 角色名称
		OperationAuths string    `db:"operation_auths"` // 权限列表
		Rank           int64     `db:"rank"`            // 排序
		ParentId       int64     `db:"parent_id"`       // 父级id
		Status         int64     `db:"status"`          // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime     time.Time `db:"create_time"`     // 创建时间
		UpdateTime     time.Time `db:"update_time"`     // 更新时间
	}
)

func NewRoleModel(conn sqlx.SqlConn, c cache.CacheConf) RoleModel {
	return &defaultRoleModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`role`",
	}
}

/*
func (m *defaultRoleModel) Insert(data *Role) (sql.Result,error) {
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, data.Id)
roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, data.Name)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, roleRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status)
	}, roleIdKey, roleNameKey)
	return ret,err
}
*/

func (m *defaultRoleModel) Insert(session sqlx.Session, data *Role) (sql.Result, error) {
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, data.Id)
	roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, data.Name)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, roleRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status)
		}
		return conn.Exec(query, data.Id, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status)
	}, roleIdKey, roleNameKey)
	return ret, err
}

func (m *defaultRoleModel) FindOne(id int64) (*Role, error) {
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, id)
	var resp Role
	err := m.QueryRow(&resp, roleIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", roleRows, m.table)
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

func (m *defaultRoleModel) FindOneByName(name string) (*Role, error) {
	roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, name)
	var resp Role
	err := m.QueryRowIndex(&resp, roleNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", roleRows, m.table)
		if err := conn.QueryRow(&resp, query, name); err != nil {
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
func (m *defaultRoleModel) Update(data *Role) error {
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, data.Id)
roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, data.Name)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, roleRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status, data.Id)
	}, roleNameKey, roleIdKey)
	return err
}
*/

func (m *defaultRoleModel) Update(session sqlx.Session, data *Role) error {
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, data.Id)
	roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, data.Name)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, roleRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status, data.Id)
		}
		return conn.Exec(query, data.Name, data.OperationAuths, data.Rank, data.ParentId, data.Status, data.Id)
	}, roleNameKey, roleIdKey)
	return err
}

func (m *defaultRoleModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	roleIdKey := fmt.Sprintf("%s%v", cacheRoleIdPrefix, id)
	roleNameKey := fmt.Sprintf("%s%v", cacheRoleNamePrefix, data.Name)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, roleNameKey, roleIdKey)
	return err
}

func (m *defaultRoleModel) SoftDelete(session sqlx.Session, data *Role) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultRoleModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheRoleIdPrefix, primary)
}

func (m *defaultRoleModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", roleRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultRoleModel) FindAll(in *tool.GetsReq) ([]*RoleItem, error) {
	resp := make([]*RoleItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, roleListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
