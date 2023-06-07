package shortener

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateShortURL(randomNumber uint64) string {
	shortURLLenght := 8
	length := len(alphabet)
	shortURL := make([]byte, shortURLLenght)
	i := 0

	for ; i < shortURLLenght; randomNumber = randomNumber / uint64(length) {
		shortURL[i] = alphabet[(randomNumber % uint64(length))]
		i++
	}

	return string(shortURL)
}
