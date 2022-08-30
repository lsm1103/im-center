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

	"im-center/common/tool"
)

var (
	offlineMsgFieldNames          = builder.RawFieldNames(&OfflineMsg{})
	offlineMsgRows                = strings.Join(offlineMsgFieldNames, ",")
	offlineMsgRowsExpectAutoSet   = strings.Join(stringx.Remove(offlineMsgFieldNames, "`create_time`", "`update_time`"), ",")
	offlineMsgRowsWithPlaceHolder = strings.Join(stringx.Remove(offlineMsgFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	offlineMsgListRows            = strings.Join(builder.RawFieldNames(&OfflineMsgItem{}), ",")

	cacheOfflineMsgIdPrefix                       = "cache:offlineMsg:id:"
	cacheOfflineMsgUserIdObjectTypeObjectIdPrefix = "cache:offlineMsg:userId:objectType:objectId:"
)

type (
	OfflineMsgModel interface {
		Insert(session sqlx.Session, data *OfflineMsg) (sql.Result, error)
		// Insert(data *OfflineMsg) (sql.Result,error)
		FindOne(id int64) (*OfflineMsg, error)
		FindAll(in *tool.GetsReq) ([]*OfflineMsgItem, error)
		FindOneByUserIdObjectTypeObjectId(userId int64, objectType int64, objectId int64) (*OfflineMsg, error)
		Update(session sqlx.Session, data *OfflineMsg) error
		// Update(data *OfflineMsg) error
		Delete(id int64) error
	}

	defaultOfflineMsgModel struct {
		sqlc.CachedConn
		table string
	}

	OfflineMsg struct {
		Id         int64     `db:"id"`           // 自增主键
		UserId     int64     `db:"user_id"`      // 用户id
		DeviceId   string    `db:"device_id"`    // 设备id
		ObjectType int64     `db:"object_type"`  // 对象类型,1:friend；2：群组
		ObjectId   int64     `db:"object_id"`    // 对象id, friendId/groupId
		LastAckSeq int64     `db:"last_ack_seq"` // 最后确认序列号
		CreateTime time.Time `db:"create_time"`  // 创建时间
		UpdateTime time.Time `db:"update_time"`  // 更新时间
	}
)

func NewOfflineMsgModel(conn sqlx.SqlConn, c cache.CacheConf) OfflineMsgModel {
	return &defaultOfflineMsgModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`offline_msg`",
	}
}

/*
func (m *defaultOfflineMsgModel) Insert(data *OfflineMsg) (sql.Result,error) {
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, data.Id)
offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, data.UserId, data.ObjectType, data.ObjectId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, offlineMsgRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq)
	}, offlineMsgIdKey, offlineMsgUserIdObjectTypeObjectIdKey)
	return ret,err
}
*/

func (m *defaultOfflineMsgModel) Insert(session sqlx.Session, data *OfflineMsg) (sql.Result, error) {
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, data.Id)
	offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, data.UserId, data.ObjectType, data.ObjectId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, offlineMsgRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq)
		}
		return conn.Exec(query, data.Id, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq)
	}, offlineMsgIdKey, offlineMsgUserIdObjectTypeObjectIdKey)
	return ret, err
}

func (m *defaultOfflineMsgModel) FindOne(id int64) (*OfflineMsg, error) {
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, id)
	var resp OfflineMsg
	err := m.QueryRow(&resp, offlineMsgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", offlineMsgRows, m.table)
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

func (m *defaultOfflineMsgModel) FindOneByUserIdObjectTypeObjectId(userId int64, objectType int64, objectId int64) (*OfflineMsg, error) {
	offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, userId, objectType, objectId)
	var resp OfflineMsg
	err := m.QueryRowIndex(&resp, offlineMsgUserIdObjectTypeObjectIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `object_type` = ? and `object_id` = ? limit 1", offlineMsgRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, objectType, objectId); err != nil {
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
func (m *defaultOfflineMsgModel) Update(data *OfflineMsg) error {
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, data.Id)
offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, data.UserId, data.ObjectType, data.ObjectId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, offlineMsgRowsWithPlaceHolder)
		return conn.Exec(query, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq, data.Id)
	}, offlineMsgIdKey, offlineMsgUserIdObjectTypeObjectIdKey)
	return err
}
*/

func (m *defaultOfflineMsgModel) Update(session sqlx.Session, data *OfflineMsg) error {
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, data.Id)
	offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, data.UserId, data.ObjectType, data.ObjectId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, offlineMsgRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq, data.Id)
		}
		return conn.Exec(query, data.UserId, data.DeviceId, data.ObjectType, data.ObjectId, data.LastAckSeq, data.Id)
	}, offlineMsgIdKey, offlineMsgUserIdObjectTypeObjectIdKey)
	return err
}

func (m *defaultOfflineMsgModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	offlineMsgIdKey := fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, id)
	offlineMsgUserIdObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v:%v", cacheOfflineMsgUserIdObjectTypeObjectIdPrefix, data.UserId, data.ObjectType, data.ObjectId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, offlineMsgUserIdObjectTypeObjectIdKey, offlineMsgIdKey)
	return err
}

func (m *defaultOfflineMsgModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheOfflineMsgIdPrefix, primary)
}

func (m *defaultOfflineMsgModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", offlineMsgRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultOfflineMsgModel) FindAll(in *tool.GetsReq) ([]*OfflineMsgItem, error) {
	resp := make([]*OfflineMsgItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, offlineMsgListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
