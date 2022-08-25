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
	seqFieldNames          = builder.RawFieldNames(&Seq{})
	seqRows                = strings.Join(seqFieldNames, ",")
	seqRowsExpectAutoSet   = strings.Join(stringx.Remove(seqFieldNames, "`create_time`", "`update_time`"), ",")
	seqRowsWithPlaceHolder = strings.Join(stringx.Remove(seqFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	seqListRows            = strings.Join(builder.RawFieldNames(&SeqItem{}), ",")

	cacheSeqIdPrefix                 = "cache:seq:id:"
	cacheSeqObjectTypeObjectIdPrefix = "cache:seq:objectType:objectId:"
)

type (
	SeqModel interface {
		Insert(session sqlx.Session, data *Seq) (sql.Result, error)
		// Insert(data *Seq) (sql.Result,error)
		FindOne(id int64) (*Seq, error)
		FindAll(in *tool.GetsReq) ([]*SeqItem, error)
		FindOneByObjectTypeObjectId(objectType int64, objectId int64) (*Seq, error)
		Update(session sqlx.Session, data *Seq) error
		// Update(data *Seq) error
		SoftDelete(session sqlx.Session, data *Seq) error
		Delete(id int64) error
	}

	defaultSeqModel struct {
		sqlc.CachedConn
		table string
	}

	Seq struct {
		Id         int64     `db:"id"`          // 自增主键
		ObjectType int64     `db:"object_type"` // 对象类型,1:用户；2：群组
		ObjectId   int64     `db:"object_id"`   // 对象id
		Seq        int64     `db:"seq"`         // 序列号
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
	}
)

func NewSeqModel(conn sqlx.SqlConn, c cache.CacheConf) SeqModel {
	return &defaultSeqModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`seq`",
	}
}

/*
func (m *defaultSeqModel) Insert(data *Seq) (sql.Result,error) {
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, data.Id)
seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, data.ObjectType, data.ObjectId)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, seqRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.ObjectType, data.ObjectId, data.Seq)
	}, seqObjectTypeObjectIdKey, seqIdKey)
	return ret,err
}
*/

func (m *defaultSeqModel) Insert(session sqlx.Session, data *Seq) (sql.Result, error) {
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, data.Id)
	seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, data.ObjectType, data.ObjectId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, seqRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.ObjectType, data.ObjectId, data.Seq)
		}
		return conn.Exec(query, data.Id, data.ObjectType, data.ObjectId, data.Seq)
	}, seqObjectTypeObjectIdKey, seqIdKey)
	return ret, err
}

func (m *defaultSeqModel) FindOne(id int64) (*Seq, error) {
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, id)
	var resp Seq
	err := m.QueryRow(&resp, seqIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", seqRows, m.table)
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

func (m *defaultSeqModel) FindOneByObjectTypeObjectId(objectType int64, objectId int64) (*Seq, error) {
	seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, objectType, objectId)
	var resp Seq
	err := m.QueryRowIndex(&resp, seqObjectTypeObjectIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `object_type` = ? and `object_id` = ? limit 1", seqRows, m.table)
		if err := conn.QueryRow(&resp, query, objectType, objectId); err != nil {
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
func (m *defaultSeqModel) Update(data *Seq) error {
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, data.Id)
seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, data.ObjectType, data.ObjectId)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, seqRowsWithPlaceHolder)
		return conn.Exec(query, data.ObjectType, data.ObjectId, data.Seq, data.Id)
	}, seqIdKey, seqObjectTypeObjectIdKey)
	return err
}
*/

func (m *defaultSeqModel) Update(session sqlx.Session, data *Seq) error {
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, data.Id)
	seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, data.ObjectType, data.ObjectId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, seqRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.ObjectType, data.ObjectId, data.Seq, data.Id)
		}
		return conn.Exec(query, data.ObjectType, data.ObjectId, data.Seq, data.Id)
	}, seqIdKey, seqObjectTypeObjectIdKey)
	return err
}

func (m *defaultSeqModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	seqIdKey := fmt.Sprintf("%s%v", cacheSeqIdPrefix, id)
	seqObjectTypeObjectIdKey := fmt.Sprintf("%s%v:%v", cacheSeqObjectTypeObjectIdPrefix, data.ObjectType, data.ObjectId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, seqObjectTypeObjectIdKey, seqIdKey)
	return err
}

func (m *defaultSeqModel) SoftDelete(session sqlx.Session, data *Seq) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultSeqModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSeqIdPrefix, primary)
}

func (m *defaultSeqModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", seqRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultSeqModel) FindAll(in *tool.GetsReq) ([]*SeqItem, error) {
	resp := make([]*SeqItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, seqListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
