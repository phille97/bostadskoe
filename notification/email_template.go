package notification

import (
	"text/template"

	"github.com/phille97/bostadskoe/provider"
)

type NewPropertiesEmailTemplateData struct {
	RecipientFirstname string
	Residences         []provider.ResidenceData
}

var NewPropertiesEmailTemplate = template.Must(template.New("NewPropertiesEmail").Parse(`
<p>Hello {{ .RecipientFirstname }}, here's some fresh apartments for you</p>
<table style="border-collapse: separate;border-spacing: 13px;">
    <thead style="font-weight: bold;">
        <tr>
            <td>Size / Rooms</td>
            <td>Location</td>
            <td>Features</td>
            <td>Rent / Month</td>
        </tr>
    </thead>
    <tbody>
    {{ range .Residences }}
        <tr>
            <td style="color: black;font-weight: bold;font-size: 25px;text-decoration: none;">{{ .LivingArea }} m<sup>2</sup> / {{ .Rooms }}</td>
            <td style="color: black;font-weight: bold;font-size: 25px;text-decoration: none;">{{ .Area }}, {{ .Municipality }}</td>
            <td>
                <table style="font-size: 9px;">
                {{ range .Features }}
                    <tr>
                        <td>{{ .Category }}.{{ .Key }}</td>
                        <td>{{ .Value }}</td>
                    </tr>
                {{ end }}
                </table>
            </td>
            <td><a href="{{ .Url }}" style="color: black;font-weight: bold;font-size: 25px;text-decoration: none;">{{ .RentPerMonth }} {{ .Currency }}</a></td>
        </tr>
    {{ end }}
    </tbody>
</table>
`))
