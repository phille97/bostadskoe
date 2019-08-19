package provider

type Residence interface {
	ID() string
	Data() ResidenceData
}

type Provider interface {
	CurrentResidences(chan Residence, chan error)
}

type ResidenceData struct {
	RentPerMonth *float64
	Currency     string
	Rooms        *float64
	LivingArea   *float64 // m²
	Features     []ResidenceFeatures
	Longitude    *float64 // WGS 84 decimal
	Latitude     *float64 // WGS 84 decimal
	Municipality *string
	Area         *string
	Url          *string
}

type ResidenceFeatures struct {
	Category string
	Key      string
	Value    interface{}
}
