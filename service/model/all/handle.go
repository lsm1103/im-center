package all

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ModelHandle interface {
	Trans(fn func(session sqlx.Session) error) error
}

type defaultModelHandle struct {
	sqlc.CachedConn
}

func NewModelHandle(conn sqlx.SqlConn, c cache.CacheConf) ModelHandle {
	return &defaultModelHandle{
		CachedConn: sqlc.NewConn(conn, c),
	}
}

func (m defaultModelHandle) Trans(fn func(session sqlx.Session) error) error {
	err := m.Transact(func(session sqlx.Session) error {
		err := fn(session)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}


// 模版
//func (m *default{{.upperStartCamelObject}}Model) Trans(fn func(session sqlx.Session)error) error  {
//	err := m.Transact(func(session sqlx.Session) error {
//		err := fn(session)
//		if err != nil{
//			return err
//		}
//		return nil
//	})
//	return err
//}