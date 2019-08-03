package senate

type Bostad struct {
	ObjektNummer string
	Url          string
	Kommun       *string
	Adress       *string
	AntalRum     *float64
	Yta          *float64
	Hyra         *float64
	Balkong      bool
	Parkering    bool
	TV           *string
	Internet     *string
	Ovrigt       *string
	LedigFran    *string
	BildUrl      *string
}

func (b Bostad) ID() string {
	return b.ObjektNummer
}
