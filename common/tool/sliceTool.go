package tool

// DelSlice 删除指定元素。
func DelStrSlice(data []string, elem string) []string {
	j := 0
	for _, v := range data {
		if v != elem {
			data[j] = v
			j++
		}
	}
	return data[:j]
}
