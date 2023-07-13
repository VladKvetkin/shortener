package entities

type URL struct {
	UUID        string `db:"id"`
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"original_url"`
}
