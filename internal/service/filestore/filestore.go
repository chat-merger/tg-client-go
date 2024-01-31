package filestore

import (
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type URI = string
type FileStore interface {
	Save(file io.Reader) (URI, error)
	Get(uri URI) (io.ReadCloser, error)
}

type Local struct {
	savePath string
}

// NewLocal создает новый экземпляр FileStoreLocal
func NewLocal(savePath string) *Local {
	return &Local{
		savePath: savePath,
	}
}

// generateURI генерирует URI на основе текущей даты и случайной строки
func (fs *Local) generateURI() URI {
	randomString := randString(3)                           // Генерируем случайную строку
	currentDate := time.Now().Format("2006-01-02-15-04-05") // Форматируем текущую дату
	return currentDate + "_" + randomString
}

// Save сохраняет файл в хранилище и возвращает сгенерированный URI
func (fs *Local) Save(file io.Reader) (URI, error) {
	uri := fs.generateURI()
	saveDir := filepath.Join(fs.savePath, filepath.Dir(uri))
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", err
	}
	saveFile, err := os.Create(filepath.Join(fs.savePath, uri))
	if err != nil {
		return "", err
	}
	defer saveFile.Close()
	_, err = io.Copy(saveFile, file)
	if err != nil {
		return "", err
	}
	return uri, nil
}

// Get возвращает файл по его URI
func (fs *Local) Get(uri URI) (io.ReadCloser, error) {
	filePath := filepath.Join(fs.savePath, uri)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Генерация случайной строки
func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
