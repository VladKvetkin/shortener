package shortener

import "math/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateID() string {
	randomNumber := rand.Uint64()

	idLength := 8
	length := len(alphabet)
	id := make([]byte, idLength)
	i := 0

	for ; i < idLength; randomNumber = randomNumber / uint64(length) {
		id[i] = alphabet[(randomNumber % uint64(length))]
		i++
	}

	return string(id)
}
