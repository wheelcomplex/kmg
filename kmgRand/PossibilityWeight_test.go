package kmgRand

func (t *Tester) TestPossibilityWeigth() {
	r := NewInt64SeedKmgRand(0)
	rander := NewPossibilityWeightRander([]float64{1e-20})
	for i := 0; i < 100; i++ {
		t.Equal(rander.ChoiceOne(r), 0)
	}
	rander = NewPossibilityWeightRander([]float64{1, 2, 3, 4})
	for i := 0; i < 100; i++ {
		ret := rander.ChoiceOne(r)
		t.Ok(ret >= 0)
		t.Ok(ret <= 3)
	}
	rander = NewPossibilityWeightRander([]float64{1, 0, 3, 0, 1})
	for i := 0; i < 100; i++ {
		ret := rander.ChoiceOne(r)
		t.Ok(ret >= 0)
		t.Ok(ret <= 4)
	}
}
