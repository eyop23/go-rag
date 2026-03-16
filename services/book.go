package services

import (
	"fmt"
	"strings"

	"github.com/eyop23/insurance-go/models"
)

func FlattenBook(b models.Book) string {
	genres := strings.Join(b.Genres, ", ")
	return fmt.Sprintf(
		"Book: %s. Genres: %s. Description: %s",
		b.Title, genres, b.Description,
	)
}
