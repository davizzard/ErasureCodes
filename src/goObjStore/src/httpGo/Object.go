package httpGo

//go:generate msgp

type Object struct {
	Name string	`msg:"name"`
	Size int	`msg:"size"`
	PartsNum int	`msg:"PartsNum"`
	ParityNum int	`msg:"ParityNum"`
}