package dependencyInjection

import "testing"
import "kmg/test"

func TestContainer(t *testing.T){
	c:=NewContainer()
	err:=c.Set("num",1,"")
	test.Assert(t,err,nil)

	num,err:=c.Get("num")
	test.Ok(t,err==nil)
	test.Assert(t,num,1)

	RequestContainer,err:=c.EnterScope(ScopeRequest)
	test.Ok(t,err==nil)

	err=c.Set("s","123",ScopeRequest)
	test.Ok(t,err==ScopeNotExistError)

	err=RequestContainer.Set("s","123",ScopeRequest)
	test.Ok(t,err==nil)

	str,err:=RequestContainer.Get("s")
	test.Assert(t,err,nil)
	test.Assert(t,str,"123")

	_,err=RequestContainer.LeaveScope()
	test.Ok(t,err==nil)
}

