package utils

import "math/rand/v2"

func MaxSize(toCut string, maxSize int) string {
	if maxSize <= 0 {
		maxSize = 100
	}
	stringSize := len(toCut)
	if stringSize <= maxSize {
		return toCut
	}
	exceeded := stringSize - maxSize
	return toCut[exceeded:]
}

func ShuffleString[T ~string](to T) T {
	inRune := []rune(to)
	rand.Shuffle(
		len(inRune), func(i, j int) {
			inRune[i], inRune[j] = inRune[j], inRune[i]
		},
	)
	return T(inRune)
}
