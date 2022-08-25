package ipTools

import (
	"testing"
)

func TestGetLocal(t *testing.T) {
	s := GetLocalIp()
	t.Log(s)
}