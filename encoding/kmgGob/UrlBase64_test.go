package kmgGob

func (s S) TestUrlBase64() {
	a := 1
	bytes, err := UrlBase64Marshal(a)
	s.Equal(err, nil)
	b := 0
	err = UrlBase64Unmarshal(bytes, &b)
	s.Equal(err, nil)
	s.Equal(b, 1)
}
