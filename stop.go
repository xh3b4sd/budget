package budget

type Stop struct{}

func (s Stop) Error() string { return "" }
