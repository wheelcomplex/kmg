package sessionStore

func (t *Tester) TestManager() {
	manager := &Manager{NewMemoryProvider()}
	session, err := manager.Load("1")
	t.Equal(err, nil)
	session.Set("A", 5)
	err = manager.Save(session)
	t.Equal(err, nil)

	sessionId := session.Id
	t.Ok(sessionId != "")

	session, err = manager.Load(sessionId)
	t.Equal(err, nil)
	value, ok := session.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 5)

	session.DeleteAndNewSession()
	_, ok = session.Get("A")
	t.Equal(ok, false)
	t.Ok(session.Id != sessionId)

	session, err = manager.Load(sessionId)
	t.Equal(err, nil)
	//pass in sessionId should not exist
	t.Ok(session.Id != sessionId)
}
