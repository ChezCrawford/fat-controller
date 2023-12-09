package conductor

type SimDriver struct{}

func NewSimDriver() DccDriver {
	return &SimDriver{}
}
func (d *SimDriver) Start() (err error) {
	return nil
}

func (d *SimDriver) SendRawCommand(rawCommand string) (err error) {
	return nil
}
