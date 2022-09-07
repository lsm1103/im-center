package database

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
	groupFieldNames          = builder.RawFieldNames(&Group{})
	groupRows                = strings.Join(groupFieldNames, ",")
	groupRowsExpectAutoSet   = strings.Join(stringx.Remove(groupFieldNames, "`create_time`", "`update_time`"), ",")
	groupRowsWithPlaceHolder = strings.Join(stringx.Remove(groupFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	groupListRows            = strings.Join(builder.RawFieldNames(&GroupItem{}), ",")

	cacheGroupIdPrefix             = "cache:group:id:"
	cacheGroupCreateUserNamePrefix = "cache:group:createUser:name:"
)

type (
	GroupModel interface {
		Insert(session sqlx.Session, data *Group) (sql.Result, error)
		// Insert(data *Group) (sql.Result,error)
		FindOne(id int64) (*Group, error)
		FindAll(in *tool.GetsReq) ([]*GroupItem, error)
		FindOneByCreateUserName(createUser int64, name string) (*Group, error)
		Update(session sqlx.Session, data *Group) error
		// Update(data *Group) error
		SoftDelete(session sqlx.Session, data *Group) error
		Delete(id int64) error
	}

	defaultGroupModel struct {
		sqlc.CachedConn
		table string
	}

	Group struct {
		Id         int64     `db:"id"`          // 自增主键
		Name       string    `db:"name"`        // 组名称
		CreateUser int64     `db:"create_user"` // 创建者id
		Ico        string    `db:"ico"`         // 组图标
		Remark     string    `db:"remark"`      // 备注
		ParentId   int64     `db:"parent_id"`   // 父级id
		GroupType  int64     `db:"group_type"`  // 类型: 1部门、2用户组、3群组、4圈子、5话题
		Rank       int64     `db:"rank"`        // 排序
		Status     int64     `db:"status"`      // 状态，-2删除，-1禁用，待审核0，启用1
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
	}
)

func NewGroupModel(conn sqlx.SqlConn, c cache.CacheConf) GroupModel {
	return &defaultGroupModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`group`",
	}
}

/*
func (m *defaultGroupModel) Insert(data *Group) (sql.Result,error) {
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, data.Id)
groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, data.CreateUser, data.Name)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.Name, data.CreateUser, data.Ico, data.Remark, data.ParentId, data.GroupType, data.Rank, data.Status)
	}, groupIdKey, groupCreateUserNameKey)
	return ret,err
}
*/

func (m *defaultGroupModel) Insert(session sqlx.Session, data *Group) (sql.Result, error) {
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, data.Id)
	groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, data.CreateUser, data.Name)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.Name, data.CreateUser, data.Ico, data.Remark, data.ParentId, data.GroupType, data.Rank, data.Status)
		}
		return conn.Exec(query, data.Id, data.Name, data.CreateUser, data.Ico, data.Remark, data.ParentId, data.GroupType, data.Rank, data.Status)
	}, groupIdKey, groupCreateUserNameKey)
	return ret, err
}

func (m *defaultGroupModel) FindOne(id int64) (*Group, error) {
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, id)
	var resp Group
	err := m.QueryRow(&resp, groupIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupRows, m.table)
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

func (m *defaultGroupModel) FindOneByCreateUserName(createUser int64, name string) (*Group, error) {
	groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, createUser, name)
	var resp Group
	err := m.QueryRowIndex(&resp, groupCreateUserNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `create_user` = ? and `name` = ? limit 1", groupRows, m.table)
		if err := conn.QueryRow(&resp, query, createUser, name); err != nil {
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
func (m *defaultGroupModel) Update(data *Group) error {
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, data.Id)
groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, data.CreateUser, data.Name)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.CreateUser, data.Ico, data.Remark, data.ParentId, data.GroupType, data.Rank, data.Status, data.Id)
	}, groupCreateUserNameKey, groupIdKey)
	return err
}
*/

func (m *defaultGroupModel) Update(session sqlx.Session, data *Group) error {
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, data.Id)
	groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, data.CreateUser, data.Name)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf( "update %s set %s where `id` = %d", m.table,
			strings.Join(tool.BuildFields(data, tool.IsEmptyValue), ", "), data.Id )
		if data.Id == 0 {
			query = fmt.Sprintf( "update %s set %s where `create_user` = %d and `name` = %d", m.table,
				strings.Join(tool.BuildFields(data, tool.IsEmptyValue), ", "), data.CreateUser, data.Name )
		}
		if session != nil {
			return session.Exec(query)
		}
		return conn.Exec(query)
	}, groupCreateUserNameKey, groupIdKey)
	return err
}

func (m *defaultGroupModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	groupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, id)
	groupCreateUserNameKey := fmt.Sprintf("%s%v:%v", cacheGroupCreateUserNamePrefix, data.CreateUser, data.Name)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, groupIdKey, groupCreateUserNameKey)
	return err
}

func (m *defaultGroupModel) SoftDelete(session sqlx.Session, data *Group) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultGroupModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGroupIdPrefix, primary)
}

func (m *defaultGroupModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultGroupModel) FindAll(in *tool.GetsReq) ([]*GroupItem, error) {
	resp := make([]*GroupItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, groupListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
