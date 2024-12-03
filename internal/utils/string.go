package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"unicode"
)

func NormalizeSpaces(s string) string {
	r := strings.Fields(s)
	return strings.Join(r, " ")
}

func Capitalize(s string) string {
	// Diviser la chaîne en mots
	s = strings.TrimSpace(s)
	words := strings.Fields(s)

	// Parcourir chaque mot
	for i, word := range words {
		// Convertir le premier caractère en majuscule et le reste en minuscule
		runes := []rune(word)
		for j := range runes {
			if j == 0 {
				runes[j] = unicode.ToUpper(runes[j])
			} else {
				runes[j] = unicode.ToLower(runes[j])
			}
		}
		words[i] = string(runes)
	}

	// Rejoindre les mots en une seule chaîne
	return strings.Join(words, " ")
}

func GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
