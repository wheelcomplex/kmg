package dependencyInjection

import "testing"
import "github.com/bronze1man/kmg/kmgTest"

func TestContainer(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	c := NewContainer()
	err := c.Set("num", 1, "")
	t.Equal(err, nil)

	num, err := c.Get("num")
	t.Ok(err == nil)
	t.Equal(num, 1)

	err = c.SetFactory("factory1", func(c *Container) (interface{}, error) {
		return 5, nil
	}, ScopeRequest)
	t.Equal(err, nil)

	RequestContainer, err := c.EnterScope(ScopeRequest)
	t.Ok(err == nil)

	err = c.Set("s", "123", ScopeRequest)
	t.Equal(err, CanNotSetNotActiveScopeByObjError)

	err = RequestContainer.Set("s", "123", ScopeRequest)
	t.Ok(err == nil)

	str, err := RequestContainer.Get("s")
	t.Equal(err, nil)
	t.Equal(str, "123")

	ret, err := RequestContainer.Get("factory1")
	t.Equal(err, nil)
	t.Equal(ret, 5)

	_, err = RequestContainer.LeaveScope()
	t.Ok(err == nil)
}
