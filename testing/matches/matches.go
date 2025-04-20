package matches

type Matcher interface {
	Matches(v any)
}

type Diff struct{}

type EqualTo struct {
	expected any
}

func (e *EqualTo) Matches(v any) *Diff {
	return nil
}
