package bostadstockholm

import (
	"strconv"
)

type Bostad struct {
	AnnonsId          int
	Stadsdel          *string
	Kommun            *string
	Vaning            *int
	AntalRum          *int
	Yta               *int
	Hyra              *int
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
		KotidFordelningQ1 int
		KotidFordelningQ2 int
	}
}

func (b Bostad) ID() string {
	return strconv.Itoa(b.AnnonsId)
}
