// Пакет entities используется для описаний сущностей базы данных.

package entities

// URL - структура, которая описывает строку таблицы url в базе данных.
type URL struct {
	UUID        string `db:"id"`
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"original_url"`
	UserID      string `db:"user_id"`
	DeletedFlag bool   `db:"is_deleted"`
}
