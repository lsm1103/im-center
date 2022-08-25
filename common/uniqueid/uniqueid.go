package uniqueid

import (
	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

var flake *sonyflake.Sonyflake

func init()  {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenId() int64 {
	id,err:= flake.NextID()
	if err != nil{
		logx.Severef("flake NextID failed with %s \n",err)
		panic(err)
	}

	return int64(id)
}

func GenUid() string {
	return strconv.FormatInt(GenId(), 10)
}