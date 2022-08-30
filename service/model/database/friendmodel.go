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
	friendFieldNames          = builder.RawFieldNames(&Friend{})
	friendRows                = strings.Join(friendFieldNames, ",")
	friendRowsExpectAutoSet   = strings.Join(stringx.Remove(friendFieldNames, "`create_time`", "`update_time`"), ",")
	friendRowsWithPlaceHolder = strings.Join(stringx.Remove(friendFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	friendListRows            = strings.Join(builder.RawFieldNames(&FriendItem{}), ",")

	cacheFriendIdPrefix                  = "cache:friend:id:"
	cacheFriendApplyUserAcceptUserPrefix = "cache:friend:applyUser:acceptUser:"
)

type (
	FriendModel interface {
		Insert(session sqlx.Session, data *Friend) (sql.Result, error)
		// Insert(data *Friend) (sql.Result,error)
		FindOne(id int64) (*Friend, error)
		FindAll(in *tool.GetsReq) ([]*FriendItem, error)
		FindOneByApplyUserAcceptUser(applyUser int64, acceptUser int64) (*Friend, error)
		Update(session sqlx.Session, data *Friend) error
		// Update(data *Friend) error
		SoftDelete(session sqlx.Session, data *Friend) error
		Delete(id int64) error
	}

	defaultFriendModel struct {
		sqlc.CachedConn
		table string
	}

	Friend struct {
		Id          int64     `db:"id"`           // 自增主键
		ApplyUser   int64     `db:"apply_user"`   // 申请用户id
		ApplyDevice string    `db:"apply_device"` // 申请设备id
		AcceptUser  int64     `db:"accept_user"`  // 接受用户id
		ApplyReason string    `db:"apply_reason"` // 申请理由
		Extra       string    `db:"extra"`        // 附加属性
		Status      int64     `db:"status"`       // 状态，-2：拉黑，-1：拒绝，0：申请中，1：同意
		CreateTime  time.Time `db:"create_time"`  // 创建时间
		UpdateTime  time.Time `db:"update_time"`  // 更新时间
	}
)

func NewFriendModel(conn sqlx.SqlConn, c cache.CacheConf) FriendModel {
	return &defaultFriendModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`friend`",
	}
}

/*
func (m *defaultFriendModel) Insert(data *Friend) (sql.Result,error) {
	friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, data.ApplyUser, data.AcceptUser)
friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, data.Id)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, friendRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status)
	}, friendApplyUserAcceptUserKey, friendIdKey)
	return ret,err
}
*/

func (m *defaultFriendModel) Insert(session sqlx.Session, data *Friend) (sql.Result, error) {
	friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, data.ApplyUser, data.AcceptUser)
	friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, data.Id)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, friendRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status)
		}
		return conn.Exec(query, data.Id, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status)
	}, friendApplyUserAcceptUserKey, friendIdKey)
	return ret, err
}

func (m *defaultFriendModel) FindOne(id int64) (*Friend, error) {
	friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, id)
	var resp Friend
	err := m.QueryRow(&resp, friendIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", friendRows, m.table)
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

func (m *defaultFriendModel) FindOneByApplyUserAcceptUser(applyUser int64, acceptUser int64) (*Friend, error) {
	friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, applyUser, acceptUser)
	var resp Friend
	err := m.QueryRowIndex(&resp, friendApplyUserAcceptUserKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `apply_user` = ? and `accept_user` = ? limit 1", friendRows, m.table)
		if err := conn.QueryRow(&resp, query, applyUser, acceptUser); err != nil {
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
func (m *defaultFriendModel) Update(data *Friend) error {
	friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, data.Id)
friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, data.ApplyUser, data.AcceptUser)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRowsWithPlaceHolder)
		return conn.Exec(query, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status, data.Id)
	}, friendIdKey, friendApplyUserAcceptUserKey)
	return err
}
*/

func (m *defaultFriendModel) Update(session sqlx.Session, data *Friend) error {
	friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, data.Id)
	friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, data.ApplyUser, data.AcceptUser)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status, data.Id)
		}
		return conn.Exec(query, data.ApplyUser, data.ApplyDevice, data.AcceptUser, data.ApplyReason, data.Extra, data.Status, data.Id)
	}, friendIdKey, friendApplyUserAcceptUserKey)
	return err
}

func (m *defaultFriendModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	friendIdKey := fmt.Sprintf("%s%v", cacheFriendIdPrefix, id)
	friendApplyUserAcceptUserKey := fmt.Sprintf("%s%v:%v", cacheFriendApplyUserAcceptUserPrefix, data.ApplyUser, data.AcceptUser)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, friendIdKey, friendApplyUserAcceptUserKey)
	return err
}

func (m *defaultFriendModel) SoftDelete(session sqlx.Session, data *Friend) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultFriendModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheFriendIdPrefix, primary)
}

func (m *defaultFriendModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", friendRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultFriendModel) FindAll(in *tool.GetsReq) ([]*FriendItem, error) {
	resp := make([]*FriendItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, friendListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
