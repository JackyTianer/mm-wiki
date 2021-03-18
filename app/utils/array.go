package utils

type array struct {
}

var Array = NewArray()

func NewArray() *array {
	return &array{}
}

func (*array) Contains(a []string, e string) bool {
	for _, a := range a {
		if a == e {
			return true
		}
	}
	return false
}
