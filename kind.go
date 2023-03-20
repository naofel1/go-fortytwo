package fortytwo

import "fmt"

// Kind defines the type for the "kind" enum field.
type Kind string

// Kind values.
const (
	KindSocial   Kind = "social"
	KindPedagogy Kind = "pedagogy"
)

func (ro Kind) String() string {
	return string(ro)
}

// KindValidator is a validator for the "Kind" field enum values. It is called by the builders before save.
func KindValidator(ro Kind) error {
	switch ro {
	case KindSocial, KindPedagogy:
		return nil
	default:
		return fmt.Errorf("invalid enum value for Kind field: %q", ro)
	}
}
