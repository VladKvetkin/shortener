package shortener

import "math/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateId() string {
	randomNumber := rand.Uint64()

	idLength := 8
	length := len(alphabet)
	uniqueId := make([]byte, idLength)
	i := 0

	for ; i < idLength; randomNumber = randomNumber / uint64(length) {
		uniqueId[i] = alphabet[(randomNumber % uint64(length))]
		i++
	}

	return string(uniqueId)
}
