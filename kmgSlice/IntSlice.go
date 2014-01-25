package kmgSlice

//remove thing at index i
func IntSliceRemoveAt(s *[]int, i int) {
	*s = append((*s)[:i], (*s)[i+1:]...)
}

//remove thing which value is v
func IntSliceRemove(s *[]int, v int) {
	thisLen := len(*s)
	for i := 0; i < thisLen; i++ {
		if (*s)[i] == v {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}
