package tool

import "testing"

func TestDelSlice(t *testing.T) {
	data := []string{"fddd", "dddd", "false","falsse"}
	s := DelStrSlice(data, "false")
	t.Log(s)
}
