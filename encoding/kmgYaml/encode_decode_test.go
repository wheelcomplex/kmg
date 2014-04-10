package kmgYaml

func (t *S) TestEncodeDecode() {
	in:=map[string]string{
		"总数  :":"总数 :",
		"1":"2",
	}
	data,err:=Marshal(in)
	t.Equal(err,nil)
	out:=map[string]string{}
	err = Unmarshal(data,&out)
	t.Equal(err,nil)
	t.Equal(in,out)
}
