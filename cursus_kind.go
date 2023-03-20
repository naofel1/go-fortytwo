package fortytwo

// CursusKind defines the type for application.
type CursusKind string

// CursusKind values.
const (
	CursusKindTest               CursusKind = "test"
	CursusKindMain               CursusKind = "main"
	CursusKindMainDeprecated     CursusKind = "main_deprecated"
	CursusKindExternal           CursusKind = "external"
	CursusKindExternalDeprecated CursusKind = "external_deprecated"
	CursusKindPiscine            CursusKind = "piscine"
	CursusKindPiscineCommunity   CursusKind = "piscine_community"
	CursusKindPiscineDeprecated  CursusKind = "piscine_deprecated"
)

// String returns the string value for CursusKind.
func (ro CursusKind) String() string {
	return string(ro)
}
