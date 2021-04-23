package cmd

type config struct {
	auth
	power
	bootDevice
	Timeout  int
	Protocol string
}

type auth struct {
	IP   string `validate:"required"`
	User string `validate:"required"`
	Pass string `validate:"required"`
}

type power struct {
	Action string `validate:"oneof=on off cycle reset"`
}

type bootDevice struct {
	Device     string `validate:"oneof=pxe disk"`
	Persistent bool
}
