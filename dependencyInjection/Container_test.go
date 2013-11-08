package dependencyInjection

import "testing"
import "kmg/test"

func TestContainer(ot *testing.T) {
	t:=test.NewTestTools(ot)
	c := NewContainer()
	err := c.Set("num", 1, "")
	t.Equal(err, nil)

	num, err := c.Get("num")
	t.Ok(err == nil)
	t.Equal(num, 1)

	RequestContainer, err := c.EnterScope(ScopeRequest)
	t.Ok(err == nil)

	err = c.Set("s", "123", ScopeRequest)
	t.Equal(err,CanNotSetNotActiveScopeByObjError)

	err = RequestContainer.Set("s", "123", ScopeRequest)
	t.Ok(err == nil)

	str, err := RequestContainer.Get("s")
	t.Equal(err, nil)
	t.Equal(str, "123")

	_, err = RequestContainer.LeaveScope()
	t.Ok(err == nil)
}
