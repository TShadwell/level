package level

type Level struct {
	UnderlyingLevel
}

func New(l UnderlyingLevel) *Level {
	return &Level{
		l,
	}
}
