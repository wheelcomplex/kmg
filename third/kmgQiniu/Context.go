package kmgQiniu

import "github.com/qiniu/api/rs"

type Context struct {
	client rs.Client
	bucket string
}

func NewContext(bucket string) *Context {
	return &Context{
		client: rs.New(nil),
		bucket: bucket,
	}
}
