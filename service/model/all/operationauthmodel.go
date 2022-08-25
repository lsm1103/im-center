package all

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"

	"im-center/common/globalkey"
	"im-center/common/tool"
)

var (
	operationAuthFieldNames          = builder.RawFieldNames(&OperationAuth{})
	operationAuthRows                = strings.Join(operationAuthFieldNames, ",")
	operationAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(operationAuthFieldNames, "`create_time`", "`update_time`"), ",")
	operationAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(operationAuthFieldNames, "`auth_identifier`", "`create_time`", "`update_time`"), "=?,") + "=?"
	operationAuthListRows            = strings.Join(builder.RawFieldNames(&OperationAuthItem{}), ",")

	cacheOperationAuthAuthIdentifierPrefix = "cache:operationAuth:authIdentifier:"
)

type (
	OperationAuthModel interface {
		Insert(session sqlx.Session, data *OperationAuth) (sql.Result, error)
		// Insert(data *OperationAuth) (sql.Result,error)
		FindOne(authIdentifier string) (*OperationAuth, error)
		FindAll(in *tool.GetsReq) ([]*OperationAuthItem, error)
		Update(session sqlx.Session, data *OperationAuth) error
		// Update(data *OperationAuth) error
		SoftDelete(session sqlx.Session, data *OperationAuth) error
		Delete(authIdentifier string) error
	}

	defaultOperationAuthModel struct {
		sqlc.CachedConn
		table string
	}

	OperationAuth struct {
		AuthIdentifier string `db:"auth_identifier"` // 权限标识（路由）
		AuthName       string `db:"auth_name"`       // 权限名称
		AuthType       string `db:"auth_type"`       // 权限类型-接口/页面/菜单/按钮
		Status         int64  `db:"status"`          // 状态，-2删除，-1禁用，待审核0，启用1
	}
)

func NewOperationAuthModel(conn sqlx.SqlConn, c cache.CacheConf) OperationAuthModel {
	return &defaultOperationAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`operation_auth`",
	}
}

/*
func (m *defaultOperationAuthModel) Insert(data *OperationAuth) (sql.Result,error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, operationAuthRowsExpectAutoSet)
    ret,err:=m.ExecNoCache(query, data.AuthIdentifier, data.AuthName, data.AuthType, data.Status)

	return ret,err
}
*/

func (m *defaultOperationAuthModel) Insert(session sqlx.Session, data *OperationAuth) (sql.Result, error) {
	operationAuthAuthIdentifierKey := fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, data.AuthIdentifier)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, operationAuthRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.AuthIdentifier, data.AuthName, data.AuthType, data.Status)
		}
		return conn.Exec(query, data.AuthIdentifier, data.AuthName, data.AuthType, data.Status)
	}, operationAuthAuthIdentifierKey)
	return ret, err
}

func (m *defaultOperationAuthModel) FindOne(authIdentifier string) (*OperationAuth, error) {
	operationAuthAuthIdentifierKey := fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, authIdentifier)
	var resp OperationAuth
	err := m.QueryRow(&resp, operationAuthAuthIdentifierKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `auth_identifier` = ? limit 1", operationAuthRows, m.table)
		return conn.QueryRow(v, query, authIdentifier)
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

/*
func (m *defaultOperationAuthModel) Update(data *OperationAuth) error {
	operationAuthAuthIdentifierKey := fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, data.AuthIdentifier)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `auth_identifier` = ?", m.table, operationAuthRowsWithPlaceHolder)
		return conn.Exec(query, data.AuthName, data.AuthType, data.Status, data.AuthIdentifier)
	}, operationAuthAuthIdentifierKey)
	return err
}
*/

func (m *defaultOperationAuthModel) Update(session sqlx.Session, data *OperationAuth) error {
	operationAuthAuthIdentifierKey := fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, data.AuthIdentifier)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `auth_identifier` = ?", m.table, operationAuthRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.AuthName, data.AuthType, data.Status, data.AuthIdentifier)
		}
		return conn.Exec(query, data.AuthName, data.AuthType, data.Status, data.AuthIdentifier)
	}, operationAuthAuthIdentifierKey)
	return err
}

func (m *defaultOperationAuthModel) Delete(authIdentifier string) error {

	operationAuthAuthIdentifierKey := fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, authIdentifier)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `auth_identifier` = ?", m.table)
		return conn.Exec(query, authIdentifier)
	}, operationAuthAuthIdentifierKey)
	return err
}

func (m *defaultOperationAuthModel) SoftDelete(session sqlx.Session, data *OperationAuth) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultOperationAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheOperationAuthAuthIdentifierPrefix, primary)
}

func (m *defaultOperationAuthModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `auth_identifier` = ? limit 1", operationAuthRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultOperationAuthModel) FindAll(in *tool.GetsReq) ([]*OperationAuthItem, error) {
	resp := make([]*OperationAuthItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, operationAuthListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
