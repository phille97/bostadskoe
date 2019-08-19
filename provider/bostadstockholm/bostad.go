package bostadstockholm

import (
	"strconv"

	"github.com/phille97/bostadskoe/provider"
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

func (b Bostad) Data() provider.ResidenceData {
	features := []provider.ResidenceFeatures{}

	features = append(features, provider.ResidenceFeatures{
		Category: "exterior",
		Key:      "balcony",
		Value:    b.Balkong,
	})

	features = append(features, provider.ResidenceFeatures{
		Category: "exterior",
		Key:      "elevator",
		Value:    b.Hiss,
	})

	if b.Student {
		features = append(features, provider.ResidenceFeatures{
			Category: "rentalconditions",
			Key:      "student",
			Value:    b.Student,
		})
	}

	if b.Ungdom {
		features = append(features, provider.ResidenceFeatures{
			Category: "rentalconditions",
			Key:      "youth",
			Value:    b.Ungdom,
		})
	}

	if b.Senior {
		features = append(features, provider.ResidenceFeatures{
			Category: "rentalconditions",
			Key:      "senior",
			Value:    b.Senior,
		})
	}

	if b.Korttid {
		features = append(features, provider.ResidenceFeatures{
			Category: "rentalconditions",
			Key:      "short-term",
			Value:    b.Korttid,
		})
	}

	return provider.ResidenceData{
		RentPerMonth: b.Hyra,
		Currency:     "SEK",
		Rooms:        b.AntalRum,
		LivingArea:   b.Yta,
		Features:     features,
		Longitude:    b.KoordinatLongitud,
		Latitude:     b.KoordinatLatitud,
		Municipality: b.Kommun,
		Area:         b.Stadsdel,
		Url:          b.Url,
	}
}
