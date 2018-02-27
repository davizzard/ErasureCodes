package httpGo
import (
	//"sync"
)
//go:generate msgp

type Account struct {
	Name string	`msg:"name"`
	Containers map[string]Container	`msg:"containers"`
}