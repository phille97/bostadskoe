package provider

type Residence interface {
	ID() string
}

type Provider interface {
	CurrentResidences() (*[]Residence, error)
}
