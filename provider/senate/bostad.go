package senate

import (
	"github.com/phille97/bostadskoe/provider"
)

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

func (b Bostad) Data() provider.ResidenceData {
	features := []provider.ResidenceFeatures{}

	features = append(features, provider.ResidenceFeatures{
		Category: "exterior",
		Key:      "balcony",
		Value:    b.Balkong,
	})

	features = append(features, provider.ResidenceFeatures{
		Category: "parking",
		Key:      "parking",
		Value:    b.Parkering,
	})

	if b.TV != nil {
		features = append(features, provider.ResidenceFeatures{
			Category: "services",
			Key:      "tv",
			Value:    b.TV,
		})
	}

	if b.Internet != nil {
		features = append(features, provider.ResidenceFeatures{
			Category: "services",
			Key:      "internet",
			Value:    b.Internet,
		})
	}

	return provider.ResidenceData{
		RentPerMonth: b.Hyra,
		Currency:     "SEK",
		Rooms:        b.AntalRum,
		Municipality: b.Kommun,
		LivingArea:   b.Yta,
		Area:         b.Adress,
		Features:     features,
		Url:          &b.Url,
	}
}
