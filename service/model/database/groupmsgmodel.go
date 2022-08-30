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
	groupMsgFieldNames          = builder.RawFieldNames(&GroupMsg{})
	groupMsgRows                = strings.Join(groupMsgFieldNames, ",")
	groupMsgRowsExpectAutoSet   = strings.Join(stringx.Remove(groupMsgFieldNames, "`create_time`", "`update_time`"), ",")
	groupMsgRowsWithPlaceHolder = strings.Join(stringx.Remove(groupMsgFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	groupMsgListRows            = strings.Join(builder.RawFieldNames(&GroupMsgItem{}), ",")

	cacheGroupMsgIdPrefix                 = "cache:groupMsg:id:"
	cacheGroupMsgSenderIdReceiverIdPrefix = "cache:groupMsg:senderId:receiverId:"
)

type (
	GroupMsgModel interface {
		Insert(session sqlx.Session, data *GroupMsg) (sql.Result, error)
		// Insert(data *GroupMsg) (sql.Result,error)
		FindOne(id int64) (*GroupMsg, error)
		FindAll(in *tool.GetsReq) ([]*GroupMsgItem, error)
		FindOneBySenderIdReceiverId(senderId int64, receiverId int64) (*GroupMsg, error)
		Update(session sqlx.Session, data *GroupMsg) error
		// Update(data *GroupMsg) error
		SoftDelete(session sqlx.Session, data *GroupMsg) error
		Delete(id int64) error
	}

	defaultGroupMsgModel struct {
		sqlc.CachedConn
		table string
	}

	GroupMsg struct {
		Id               int64     `db:"id"`                 // 自增主键(消息序列号,每个单聊都维护一个消息序列号)
		SenderType       int64     `db:"sender_type"`        // 发送者类型：1群内，2转发
		SenderId         int64     `db:"sender_id"`          // 发送者id
		SenderDeviceId   string    `db:"sender_device_id"`   // 发送设备id
		ReceiverId       int64     `db:"receiver_id"`        // 接收者id(group_id)
		ReceiverDeviceId string    `db:"receiver_device_id"` // 接收设备id：多个设备id"，"隔开，*表示全部设备
		AtUserIds        string    `db:"at_user_ids"`        // 需要@的用户id列表，多个用户用@隔开
		MsgType          int64     `db:"msg_type"`           // 消息类型：0文本、1图文、2语音、3视频、地址、4链接
		Content          string    `db:"content"`            // 消息内容
		ParentId         int64     `db:"parent_id"`          // 父级id，引用功能
		SendTime         time.Time `db:"send_time"`          // 消息发送时间
		Status           int64     `db:"status"`             // 消息状态，-3接收者删除，-2发送者删除，-1撤回，0未处理，1已读
		CreateTime       time.Time `db:"create_time"`        // 创建时间
		UpdateTime       time.Time `db:"update_time"`        // 更新时间
	}
)

func NewGroupMsgModel(conn sqlx.SqlConn, c cache.CacheConf) GroupMsgModel {
	return &defaultGroupMsgModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`group_msg`",
	}
}

/*
func (m *defaultGroupMsgModel) Insert(data *GroupMsg) (sql.Result,error) {
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, data.Id)
groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupMsgRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
	}, groupMsgIdKey, groupMsgSenderIdReceiverIdKey)
	return ret,err
}
*/

func (m *defaultGroupMsgModel) Insert(session sqlx.Session, data *GroupMsg) (sql.Result, error) {
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, data.Id)
	groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupMsgRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
		}
		return conn.Exec(query, data.Id, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status)
	}, groupMsgIdKey, groupMsgSenderIdReceiverIdKey)
	return ret, err
}

func (m *defaultGroupMsgModel) FindOne(id int64) (*GroupMsg, error) {
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, id)
	var resp GroupMsg
	err := m.QueryRow(&resp, groupMsgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMsgRows, m.table)
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

func (m *defaultGroupMsgModel) FindOneBySenderIdReceiverId(senderId int64, receiverId int64) (*GroupMsg, error) {
	groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, senderId, receiverId)
	var resp GroupMsg
	err := m.QueryRowIndex(&resp, groupMsgSenderIdReceiverIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `sender_id` = ? and `receiver_id` = ? limit 1", groupMsgRows, m.table)
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
func (m *defaultGroupMsgModel) Update(data *GroupMsg) error {
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, data.Id)
groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupMsgRowsWithPlaceHolder)
		return conn.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
	}, groupMsgIdKey, groupMsgSenderIdReceiverIdKey)
	return err
}
*/

func (m *defaultGroupMsgModel) Update(session sqlx.Session, data *GroupMsg) error {
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, data.Id)
	groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupMsgRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
		}
		return conn.Exec(query, data.SenderType, data.SenderId, data.SenderDeviceId, data.ReceiverId, data.ReceiverDeviceId, data.AtUserIds, data.MsgType, data.Content, data.ParentId, data.SendTime, data.Status, data.Id)
	}, groupMsgIdKey, groupMsgSenderIdReceiverIdKey)
	return err
}

func (m *defaultGroupMsgModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	groupMsgIdKey := fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, id)
	groupMsgSenderIdReceiverIdKey := fmt.Sprintf("%s%v:%v", cacheGroupMsgSenderIdReceiverIdPrefix, data.SenderId, data.ReceiverId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, groupMsgIdKey, groupMsgSenderIdReceiverIdKey)
	return err
}

func (m *defaultGroupMsgModel) SoftDelete(session sqlx.Session, data *GroupMsg) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultGroupMsgModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGroupMsgIdPrefix, primary)
}

func (m *defaultGroupMsgModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", groupMsgRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultGroupMsgModel) FindAll(in *tool.GetsReq) ([]*GroupMsgItem, error) {
	resp := make([]*GroupMsgItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, groupMsgListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
