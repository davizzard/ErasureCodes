package httpGo

//go:generate msgp
type Container struct {
	Name string	`msg:"name"`
	Objs map[string]Object	`msg:"objs"`
	Policy string	`msg:"policy"`
	}