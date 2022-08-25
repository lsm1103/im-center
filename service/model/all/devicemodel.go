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
	deviceFieldNames          = builder.RawFieldNames(&Device{})
	deviceRows                = strings.Join(deviceFieldNames, ",")
	deviceRowsExpectAutoSet   = strings.Join(stringx.Remove(deviceFieldNames, "`create_time`", "`update_time`"), ",")
	deviceRowsWithPlaceHolder = strings.Join(stringx.Remove(deviceFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	deviceListRows            = strings.Join(builder.RawFieldNames(&DeviceItem{}), ",")

	cacheDeviceIdPrefix                = "cache:device:id:"
	cacheDeviceUserAgentClientIpPrefix = "cache:device:userAgent:clientIp:"
)

type (
	DeviceModel interface {
		Insert(session sqlx.Session, data *Device) (sql.Result, error)
		// Insert(data *Device) (sql.Result,error)
		FindOne(id int64) (*Device, error)
		FindAll(in *tool.GetsReq) ([]*DeviceItem, error)
		FindOneByUserAgentClientIp(userAgent string, clientIp string) (*Device, error)
		Update(session sqlx.Session, data *Device) error
		// Update(data *Device) error
		SoftDelete(session sqlx.Session, data *Device) error
		Delete(id int64) error
	}

	defaultDeviceModel struct {
		sqlc.CachedConn
		table string
	}

	Device struct {
		Id            int64     `db:"id"`             // 自增主键
		UserId        int64     `db:"user_id"`        // 用户id
		UserAgent     string    `db:"user_agent"`     // 用户标示
		MachineType   int64     `db:"machine_type"`   // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
		Brand         string    `db:"brand"`          // 设备厂商
		UnitType      string    `db:"unit_type"`      // 设备型号
		SystemVersion string    `db:"system_version"` // 系统版本
		Browser       string    `db:"browser"`        // 浏览器
		Language      string    `db:"language"`       // 语言
		NetType       string    `db:"net_type"`       // 网络类型
		SdkVersion    string    `db:"sdk_version"`    // sdk版本
		ConnIp        string    `db:"conn_ip"`        // 连接层服务器ip
		ClientIp      string    `db:"client_ip"`      // 客户端ip
		ClientAddr    string    `db:"client_addr"`    // 客户端地址,浙江省杭州市西湖区浙江大学国家大学科技园
		Status        int64     `db:"status"`         // 设备状态，-2删除，-1禁用，0待审核，1离线，2在线
		CreateTime    time.Time `db:"create_time"`    // 创建时间
		UpdateTime    time.Time `db:"update_time"`    // 更新时间
	}
)

func NewDeviceModel(conn sqlx.SqlConn, c cache.CacheConf) DeviceModel {
	return &defaultDeviceModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`device`",
	}
}

/*
func (m *defaultDeviceModel) Insert(data *Device) (sql.Result,error) {
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, data.Id)
deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, data.UserAgent, data.ClientIp)
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, deviceRowsExpectAutoSet)
		return conn.Exec(query, data.Id, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status)
	}, deviceIdKey, deviceUserAgentClientIpKey)
	return ret,err
}
*/

func (m *defaultDeviceModel) Insert(session sqlx.Session, data *Device) (sql.Result, error) {
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, data.Id)
	deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, data.UserAgent, data.ClientIp)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, deviceRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status)
		}
		return conn.Exec(query, data.Id, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status)
	}, deviceIdKey, deviceUserAgentClientIpKey)
	return ret, err
}

func (m *defaultDeviceModel) FindOne(id int64) (*Device, error) {
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, id)
	var resp Device
	err := m.QueryRow(&resp, deviceIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceRows, m.table)
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

func (m *defaultDeviceModel) FindOneByUserAgentClientIp(userAgent string, clientIp string) (*Device, error) {
	deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, userAgent, clientIp)
	var resp Device
	err := m.QueryRowIndex(&resp, deviceUserAgentClientIpKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_agent` = ? and `client_ip` = ? limit 1", deviceRows, m.table)
		if err := conn.QueryRow(&resp, query, userAgent, clientIp); err != nil {
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
func (m *defaultDeviceModel) Update(data *Device) error {
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, data.Id)
deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, data.UserAgent, data.ClientIp)
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, deviceRowsWithPlaceHolder)
		return conn.Exec(query, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status, data.Id)
	}, deviceIdKey, deviceUserAgentClientIpKey)
	return err
}
*/

func (m *defaultDeviceModel) Update(session sqlx.Session, data *Device) error {
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, data.Id)
	deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, data.UserAgent, data.ClientIp)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, deviceRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status, data.Id)
		}
		return conn.Exec(query, data.UserId, data.UserAgent, data.MachineType, data.Brand, data.UnitType, data.SystemVersion, data.Browser, data.Language, data.NetType, data.SdkVersion, data.ConnIp, data.ClientIp, data.ClientAddr, data.Status, data.Id)
	}, deviceIdKey, deviceUserAgentClientIpKey)
	return err
}

func (m *defaultDeviceModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	deviceIdKey := fmt.Sprintf("%s%v", cacheDeviceIdPrefix, id)
	deviceUserAgentClientIpKey := fmt.Sprintf("%s%v:%v", cacheDeviceUserAgentClientIpPrefix, data.UserAgent, data.ClientIp)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, deviceIdKey, deviceUserAgentClientIpKey)
	return err
}

func (m *defaultDeviceModel) SoftDelete(session sqlx.Session, data *Device) error {
	data.Status = globalkey.UserDel
	return m.Update(session, data)
}

func (m *defaultDeviceModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDeviceIdPrefix, primary)
}

func (m *defaultDeviceModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", deviceRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultDeviceModel) FindAll(in *tool.GetsReq) ([]*DeviceItem, error) {
	resp := make([]*DeviceItem, 0)
	queryStr := tool.NewModelTool().BuildQuery(in, deviceListRows, m.table)
	err := m.CachedConn.QueryRowsNoCache(&resp, queryStr)
	return resp, err
}
