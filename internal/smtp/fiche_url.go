package smtp

import "fmt"

// FicheSource identifies one of the four fiche types exposed by the API.
// Values match the API path segments (mobiliers_images, monuments_lieux, …)
// and the frontend router `source` parameter.
type FicheSource string

const (
	SourceMobiliersImages    FicheSource = "mobiliers_images"
	SourceMonumentsLieux     FicheSource = "monuments_lieux"
	SourcePersonnesMorales   FicheSource = "personnes_morales"
	SourcePersonnesPhysiques FicheSource = "personnes_physiques"
)

const frontendBaseURL = "https://saintjacquesinfo.eu"

// FicheURL builds the public URL to view a fiche in the frontend router.
func FicheURL(source FicheSource, id int32) string {
	return fmt.Sprintf("%s/fiches/%s/%d", frontendBaseURL, source, id)
}

// AdminValidationURL builds the admin URL used to review a pending fiche.
func AdminValidationURL(source FicheSource, id int32) string {
	return fmt.Sprintf("%s/admin/validation-fiches/%s/%d", frontendBaseURL, source, id)
}
