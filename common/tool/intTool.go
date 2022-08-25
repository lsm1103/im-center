package tool

func Abs(v int64) int64 {
	if v == 0 {
		return 0
	}else if v < 0 {
		return -v
	}
	return v
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}