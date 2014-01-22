package sessionStore

func (t *Tester) TestManager() {
	manager := &Manager{NewMemoryProvider()}
	session, err := manager.Load("1")
	t.Equal(err, nil)
	session.Set("A", 5)
	err = manager.Save(session)
	t.Equal(err, nil)

	session, err = manager.Load(session.Id)
	t.Equal(err, nil)
	value, ok := session.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 5)
}
