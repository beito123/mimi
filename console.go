package mimi

type Console struct {
	Cmder *Cmder
}

func (con *Console) Line() (string, bool) {
	return "", false
}

func (con *Console) SendCommand(cmd string) {
	//
}
