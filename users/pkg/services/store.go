package services

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(userID int32, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

type ImageInfo struct {
	UserID int32
	Type   string
	Path   string
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

func (store *DiskImageStore) Save(userID int32, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Cannot create image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Cannot write image to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		UserID: userID,
		Type:   imageType,
		Path:   imagePath,
	}

	// fileName := fmt.Sprintf("%s.%s", imageID.String(), imageType)
	return imageID.String(), nil
}
