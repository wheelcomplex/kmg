package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	"time"
)

func TestFormat(ot *testing.T) {
	tc := kmgTest.NewTestTools(ot)
	t, err := time.Parse(AppleJsonFormat, "2014-04-16 18:26:18 Etc/GMT")
	tc.Equal(err, nil)
	tc.Ok(t.Equal(MustFromMysqlFormat("2014-04-16 18:26:18")))
}
