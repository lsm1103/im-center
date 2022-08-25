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
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`create_time`", "`update_time`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	userListRows            = strings.Join(builder.RawFieldNames(&UserItem{}), ",")

	cacheUserIdPrefix = "cache:user:id:"
)

type (
	UserModel interface {
		Insert(session sqlx.Session, data *User) (sql.Result, error)
		// Insert(data *User) (sql.Result,error)
		FindOne(id int64) (*User, error)
		FindAll(in *tool.GetsReq) ([]*UserItem, error)
		Update(session sqlx.Session, data *User) error
		// Update(data *User) error
		SoftDelete(session sqlx.Session, data *User) error
		Delete(id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id             int64     `db:"id"`              // 自增主键
		Nickname       string    `db:"nickname"`        // 昵称
		RealName       string    `db:"realName"`        // 真实姓名
		Password       string    `db:"password"`        // 密码
		LoginSalt      string    `db:"login_salt"`      // 密码加密的盐（md5计算id+注册时间）
		RegisterDevice string    `db:"register_device"` // 注册设备信息
		Sex            int64     `db:"sex"`             // 性别，0:未知；1:男；2:女
		Ico            string    `db:"ico"`             // 用户头像
		Status         int64     `db:"status"`          // 用户状态，-2:删除；-1:冻结；0：待审核；1：正常；2：第三方直接注册登入；9：超管
		CreateTime     time.Time `db:"create_time"`     // 创建/注册时间
		UpdateTime     time.Time `db:"update_time"`     // 更新时间
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

/*
func (m *defaultUserModel) Insert(data *User) (sql.Result,error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
    ret,err:=m.ExecNoCache(query, data.Id, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status)

	return ret,err
}
*/

func (m *defaultUserModel) Insert(session sqlx.Session, data *User) (sql.Result, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status)
		}
		return conn.Exec(query, data.Id, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status)
	}, userIdKey)
	return ret, err
}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
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
func (m *defaultUserModel) Update(data *User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		return conn.Exec(query, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status, data.Id)
	}, userIdKey)
	return err
}
*/

func (m *defaultUserModel) Update(session sqlx.Session, data *User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status, data.Id)
		}
		return conn.Exec(query, data.Nickname, data.RealName, data.Password, data.LoginSalt, data.RegisterDevice, data.Sex, data.Ico, data.Status, data.Id)
	}, userIdKey)
	return err
}

func (m *defaultUserModel) Delete(id int64) error {

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userIdKey)
	return err
}

func (m *defaultUserModel) SoftDelete(session sqlx.Session, data *User) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultUserModel) FindAll(in *tool.GetsReq) ([]*UserItem, error) {
	resp := make([]*UserItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, userListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
