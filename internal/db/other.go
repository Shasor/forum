package db

import (
	"encoding/base64"
	_ "image/gif"  // Import pour le décodage GIF
	_ "image/jpeg" // Import pour le décodage JPEG
	_ "image/png"  // Import pour le décodage PNG
	"os"
)

func ImageToBase64(imagePath string) (string, error) {
	// Ouvrir le fichier image
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Obtenir les informations du fichier
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	// Lire le contenu du fichier
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Encoder en base64
	base64String := base64.StdEncoding.EncodeToString(buffer)

	return base64String, nil
}
