package bostadstockholm

import (
	"strconv"
)

type Bostad struct {
	AnnonsId          int
	Stadsdel          *string
	Kommun            *string
	Vaning            *float64
	AntalRum          *float64
	Yta               *float64
	Hyra              *float64
	AnnonseradTill    *string
	AnnonseradFran    *string
	KoordinatLongitud *float64
	KoordinatLatitud  *float64
	Url               *string
	Antal             *int
	Balkong           bool
	Hiss              bool
	Nyproduktion      bool
	Ungdom            bool
	Student           bool
	Senior            bool
	Korttid           bool
	Vanlig            bool
	Bostadssnabben    bool
	Ko                *string
	KoNamn            *string
	Lagenhetstyp      *string
	HarAnmaltIntresse bool
	KanAnmalaIntresse bool
	HarBraChans       bool
	HarInternko       bool
	Internko          bool
	Externko          bool
	Omraden           []struct {
		Id       int
		PlatsTyp int
	}
	ArInloggad                bool
	LiknandeLagenhetStatistik *struct {
		KotidFordelningQ1 float64
		KotidFordelningQ2 float64
	}
}

func (b Bostad) ID() string {
	return strconv.Itoa(b.AnnonsId)
}
