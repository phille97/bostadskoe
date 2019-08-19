package provider

type Residence interface {
	ID() string
}

type Provider interface {
	CurrentResidences(chan Residence, chan error)
}
