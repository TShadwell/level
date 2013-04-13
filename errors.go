package level

type Error uint8

const (
	Already_Open Error = iota
	Not_Opened
)

func (e Error) Error() (o string) {
	switch e {
	case Already_Open:
		o = "Database was already open."
	case Not_Opened:
		o = "Database has not yet been opened."
	}
	return
}
