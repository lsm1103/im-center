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
	toDoFieldNames          = builder.RawFieldNames(&ToDo{})
	toDoRows                = strings.Join(toDoFieldNames, ",")
	toDoRowsExpectAutoSet   = strings.Join(stringx.Remove(toDoFieldNames, "`create_time`", "`update_time`"), ",")
	toDoRowsWithPlaceHolder = strings.Join(stringx.Remove(toDoFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	toDoListRows            = strings.Join(builder.RawFieldNames(&ToDoItem{}), ",")

	cacheToDoIdPrefix = "cache:toDo:id:"
)

type (
	ToDoModel interface {
		Insert(session sqlx.Session, data *ToDo) (sql.Result, error)
		// Insert(data *ToDo) (sql.Result,error)
		FindOne(id int64) (*ToDo, error)
		FindAll(in *tool.GetsReq) ([]*ToDoItem, error)
		Update(session sqlx.Session, data *ToDo) error
		// Update(data *ToDo) error
		SoftDelete(session sqlx.Session, data *ToDo) error
		Delete(id int64) error
	}

	defaultToDoModel struct {
		sqlc.CachedConn
		table string
	}

	ToDo struct {
		Id            int64        `db:"id"`             // 自增主键
		Content       string       `db:"content"`        // 待办内容
		CreateUser    int64        `db:"create_user"`    // 创建人
		ExecuteUsers  string       `db:"execute_users"`  // 执行者，可多个
		JoinUsers     string       `db:"join_users"`     // 参与者，可多个
		TodoType      int64        `db:"todo_type"`      // 待办类型，审批(他人指定)、个人(自己写的)、群公告(他人指定)、项目任务(他人指定)
		Status        int64        `db:"status"`         // 待办状态；-2：暂停，-1：删除，0：待启动，1：进行中，2：完成
		StartTime     sql.NullTime `db:"start_time"`     // 开始时间
		Deadline      sql.NullTime `db:"deadline"`       // 截止时间
		Remark        string       `db:"remark"`         // 备注,可以分割线的形式添加进度描述
		Priority      int64        `db:"priority"`       // 优先级；
		ParentId      int64        `db:"parent_id"`      // 父级id
		AttachmentUrl string       `db:"attachment_url"` // 附件(文件中心的下载地址)，可多个
		CreateTime    time.Time    `db:"create_time"`    // 创建时间
		UpdateTime    time.Time    `db:"update_time"`    // 更新时间
	}
)

func NewToDoModel(conn sqlx.SqlConn, c cache.CacheConf) ToDoModel {
	return &defaultToDoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`to_do`",
	}
}

/*
func (m *defaultToDoModel) Insert(data *ToDo) (sql.Result,error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, toDoRowsExpectAutoSet)
    ret,err:=m.ExecNoCache(query, data.Id, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl)

	return ret,err
}
*/

func (m *defaultToDoModel) Insert(session sqlx.Session, data *ToDo) (sql.Result, error) {
	toDoIdKey := fmt.Sprintf("%s%v", cacheToDoIdPrefix, data.Id)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, toDoRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl)
		}
		return conn.Exec(query, data.Id, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl)
	}, toDoIdKey)
	return ret, err
}

func (m *defaultToDoModel) FindOne(id int64) (*ToDo, error) {
	toDoIdKey := fmt.Sprintf("%s%v", cacheToDoIdPrefix, id)
	var resp ToDo
	err := m.QueryRow(&resp, toDoIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", toDoRows, m.table)
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

/*
func (m *defaultToDoModel) Update(data *ToDo) error {
	toDoIdKey := fmt.Sprintf("%s%v", cacheToDoIdPrefix, data.Id)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, toDoRowsWithPlaceHolder)
		return conn.Exec(query, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl, data.Id)
	}, toDoIdKey)
	return err
}
*/

func (m *defaultToDoModel) Update(session sqlx.Session, data *ToDo) error {
	toDoIdKey := fmt.Sprintf("%s%v", cacheToDoIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, toDoRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl, data.Id)
		}
		return conn.Exec(query, data.Content, data.CreateUser, data.ExecuteUsers, data.JoinUsers, data.TodoType, data.Status, data.StartTime, data.Deadline, data.Remark, data.Priority, data.ParentId, data.AttachmentUrl, data.Id)
	}, toDoIdKey)
	return err
}

func (m *defaultToDoModel) Delete(id int64) error {

	toDoIdKey := fmt.Sprintf("%s%v", cacheToDoIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, toDoIdKey)
	return err
}

func (m *defaultToDoModel) SoftDelete(session sqlx.Session, data *ToDo) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultToDoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheToDoIdPrefix, primary)
}

func (m *defaultToDoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", toDoRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultToDoModel) FindAll(in *tool.GetsReq) ([]*ToDoItem, error) {
	resp := make([]*ToDoItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, toDoListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
