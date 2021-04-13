package labels

//Label represents a returned Resource Type label from the API
type Label interface {
	MarshalJSON() ([]byte, error)
	String() string
}
