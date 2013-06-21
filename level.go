package level

/*
	Type Level represents a LevelDB implementation.
*/
type Level struct {
	UnderlyingLevel
}

/*
	Function New returns a *Level corresponding to the passed UnderlyingLevel.
*/
func New(l UnderlyingLevel) *Level {
	return &Level{
		l,
	}
}
