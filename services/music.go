package services

import (
	"fmt"
	"strings"

	"github.com/eyop23/insurance-go/models"
)

func FlattenArtist(a models.MusicArtist) string {
	genres := strings.Join(a.Genre, ", ")
	instruments := strings.Join(a.Instruments, ", ")
	songs := strings.Join(a.FamousSongs, ", ")
	awards := strings.Join(a.Awards, ", ")

	albumText := "No albums recorded"
	if len(a.Albums) > 0 {
		parts := make([]string, len(a.Albums))
		for i, al := range a.Albums {
			parts[i] = fmt.Sprintf("%s (%d, %s)", al.Title, al.Year, al.Label)
		}
		albumText = strings.Join(parts, "; ")
	}

	deathInfo := ""
	if a.DeathYear > 0 {
		deathInfo = fmt.Sprintf(" Died: %d.", a.DeathYear)
	}

	return fmt.Sprintf(
		"Artist: %s. ID: %s. Born: %d, %s.%s Genre: %s. Instruments: %s. Era: %s. Bio: %s Albums: %s. Famous Songs: %s. Awards: %s. Influence: %s.",
		a.Name, a.ID, a.BirthYear, a.Origin, deathInfo, genres, instruments, a.Era, a.Bio, albumText, songs, awards, a.Influence,
	)
}
