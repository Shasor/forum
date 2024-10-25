package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif" // Import pour le décodage GIF
	"image/jpeg"
	_ "image/png" // Import pour le décodage PNG
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
)

func OpenLocalImage(filePath string) (multipart.File, *multipart.FileHeader, error) {
	// Ouvrir le fichier local
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Obtenir les informations du fichier
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, err
	}

	// Créer un multipart.FileHeader
	header := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     fileInfo.Size(),
	}

	// Créer un multipart.File à partir du fichier ouvert
	multipartFile := multipart.File(file)

	return multipartFile, header, nil
}

func ImageToBase64(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Reset the file pointer to the beginning of the file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error resetting file pointer: %v", err)
	}

	// Check the file extension to see if it's a GIF
	ext := filepath.Ext(header.Filename)
	if ext == ".gif" {
		// For GIFs, just read the entire file and base64 encode it
		data, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("error reading GIF file: %v", err)
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}

	// For non-GIF images (JPG, PNG), proceed with decoding and resizing
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Resize the image (optional, depends on your requirements)
	resizedImg := resizeImage(img)

	// Encode resized image to JPEG format
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return "", fmt.Errorf("error encoding image to JPEG: %v", err)
	}

	// Base64-encode the JPEG image
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Resize img to 1920x1080
func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	maxWidth := 400
	height := bounds.Dy()
	maxHeight := 400

	// Calculez le ratio pour redimensionner
	ratioW := float64(maxWidth) / float64(width)
	ratioH := float64(maxHeight) / float64(height)
	ratio := math.Min(ratioW, ratioH)

	if ratio >= 1 {
		return img
	}

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	// Créez une nouvelle image redimensionnée
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Redimensionnez l'image manuellement
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}
