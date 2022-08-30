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
	singleMsgFieldNames          = builder.RawFieldNames(&SingleMsg{})
	singleMsgRows                = strings.Join(singleMsgFieldNames, ",")
	singleMsgRowsExpectAutoSet   = strings.Join(stringx.Remove(singleMsgFieldNames, "`create_time`", "`update_time`"), ",")
	singleMsgRowsWithPlaceHolder = strings.Join(stringx.Remove(singleMsgFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	singleMsgListRows            = strings.Join(builder.RawFieldNames(&SingleMsgItem{}), ",")

	cacheSingleMsgIdPrefix                 = "cache:singleMsg:id:"
	cacheSingleMsgSenderIdReceiverIdPrefix = "cache:singleMsg:senderId:receiverId:"
)

type (
	SingleMsgModel interface {
		Insert(session sqlx.Session, data *SingleMsg) (sql.Result, error)
		// Insert(data *SingleMsg) (sql.Result,error)
		FindOne(id int64) (*SingleMsg, error)
		FindAll(in *tool.GetsReq) ([]*SingleMsgItem, error)
		FindOneBySenderIdReceiverId(senderId int64, receiverId int64) (*SingleMsg, error)
		Update(session sqlx.Session, data *SingleMsg) error
		// Update(data *SingleMsg) error
		SoftDelete(session sqlx.Session, data *SingleMsg) error
		Delete(id int64) error
	}

	defaultSingleMsgModel struct {
		sqlc.CachedConn
		table string
	}

	SingleMsg struct {
		Id               int64     `db:"id"`                 // 自增主键(消息序列号,每个单聊都维护一个消息序列号)
		SenderType       int64     `db:"sender_type"`        // 发送者类型：1朋友，2打招呼，3转发
		SenderId         int64     `db:"sender_id"`          // 发送者id
		SenderDeviceId   string    `db:"sender_device_id"`   // 发送设备id
		ReceiverId       int64     `db:"receiver_id"`        // 接收者id
		ReceiverDeviceId string    `db:"receiver_device_id"` // 接收设备id：多个设备id"，"隔开，*表示全部设备
		MsgType          int64     `db:"msg_type"`           // 消息类型：0文本、1图文、2语音、3视频、4链接
		Content          string    `db:"content"`            // 消息内容
		ParentId         int64     `db:"parent_id"`          // 父级id，引用功能
		SendTime         time.Time `db:"send_time"`          // 消息发送时间
		Status           int64     `db:"status"`             // 消息状态：-1撤回，0未处理，1已读
		CreateTime       time.Time `db:"create_time"`        // 创建时间
		UpdateTime       time.Time `db:"update_time"`        // 更新时间
	}
)

func NewSingleMsgModel(conn sqlx.SqlConn, c cache.CacheConf) SingleMsgModel {
	return &defaultSingleMsgModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`single_msg`",
	}
}

/*
func (m *defaultSingleMsgModel) Insert(data *SingleMsg) (sql.Result,error) {
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, data.Id)
singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, singleMsgRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
	}, singleMsgIdKey, singleMsgSenderIdReceiverIdKey)
	return ret,err
}
*/

func (m *defaultSingleMsgModel) Insert(session sqlx.Session, data *SingleMsg) (sql.Result, error) {
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, data.Id)
	singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, singleMsgRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
		}
		return conn.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
	}, singleMsgIdKey, singleMsgSenderIdReceiverIdKey)
	return ret, err
}

func (m *defaultSingleMsgModel) FindOne(id int64) (*SingleMsg, error) {
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, id)
	var resp SingleMsg
	err := m.QueryRow(&resp, singleMsgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", singleMsgRows, m.table)
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

func (m *defaultSingleMsgModel) FindOneBySenderIdReceiverId(senderId int64, receiverId int64) (*SingleMsg, error) {
	singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, senderId, receiverId)
	var resp SingleMsg
	err := m.QueryRowIndex(&resp, singleMsgSenderIdReceiverIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `sender_id` = ? and `receiver_id` = ? limit 1", singleMsgRows, m.table)
		if err := conn.QueryRow(&resp, query, senderId, receiverId); err != nil {
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
func (m *defaultSingleMsgModel) Update(data *SingleMsg) error {
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, data.Id)
singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, singleMsgRowsWithPlaceHolder)
		return conn.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
	}, singleMsgIdKey, singleMsgSenderIdReceiverIdKey)
	return err
}
*/

func (m *defaultSingleMsgModel) Update(session sqlx.Session, data *SingleMsg) error {
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, data.Id)
	singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, singleMsgRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
		}
		return conn.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
	}, singleMsgIdKey, singleMsgSenderIdReceiverIdKey)
	return err
}

func (m *defaultSingleMsgModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	singleMsgIdKey := fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, id)
	singleMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheSingleMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, singleMsgIdKey, singleMsgSenderIdReceiverIdKey)
	return err
}

func (m *defaultSingleMsgModel) SoftDelete(session sqlx.Session, data *SingleMsg) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultSingleMsgModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSingleMsgIdPrefix, primary)
}

func (m *defaultSingleMsgModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", singleMsgRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultSingleMsgModel) FindAll(in *tool.GetsReq) ([]*SingleMsgItem, error) {
	resp := make([]*SingleMsgItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, singleMsgListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
