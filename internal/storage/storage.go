package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

// FileReader — интерфейс для абстракции чтения/записи файлов (DI).
type FileReader interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// OSFileReader — реализация FileReader поверх os.
type OSFileReader struct{}

// ReadFile — чтение файла через os.ReadFile.
func (oSFileReader OSFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile — запись файла через os.WriteFile.
func (oSFileReader OSFileReader) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// Storage хранит данные в памяти и управляет чтением/записью JSON-файла.
type Storage struct {
	mu          sync.Mutex
	Data        map[int]map[string]string
	LastLinkNum int
	filePath    string
	reader      FileReader
}

// storageFileData — DTO для сериализации данных в JSON.
type storageFileData struct {
	LastLinkNum int                       `json:"last_link_num"`
	Data        map[int]map[string]string `json:"data"`
}

// NewStorage создаёт объект хранилища и сразу загружает данные с диска.
func NewStorage(filePath string, reader FileReader) (*Storage, error) {
	if reader == nil {
		reader = OSFileReader{}
	}

	strg := &Storage{
		Data:        make(map[int]map[string]string),
		LastLinkNum: 0,
		filePath:    filePath,
		reader:      reader,
	}

	if err := strg.LoadFromDisk(); err != nil {
		log.Printf("storage file error: %v — starting with empty storage", err)
	}

	return strg, nil
}

// LoadFromDisk загружает данные хранилища из JSON-файла.
func (strg *Storage) LoadFromDisk() error {
	// Блокируем доступ, чтобы избежать гонок данных
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData, err := strg.reader.ReadFile(strg.filePath)
	if err != nil {
		// Файл отсутствует — инициализируем пустое состояние
		if errors.Is(err, os.ErrNotExist) {
			strg.Data = make(map[int]map[string]string)
			strg.LastLinkNum = 0
			return nil
		}
		return err
	}

	// Превращаем JSON-файл в структуру Go, чтобы можно было извлечь LastLinkNum и Data
	var parsed storageFileData
	if err := json.Unmarshal(fileData, &parsed); err != nil {
		return err
	}

	if parsed.Data == nil {
		parsed.Data = make(map[int]map[string]string)
	}

	strg.Data = parsed.Data
	strg.LastLinkNum = parsed.LastLinkNum

	return nil
}

// SaveToDisk сохраняет текущее состояние хранилища в JSON-файл.
func (strg *Storage) SaveToDisk() error {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData := &storageFileData{
		LastLinkNum: strg.LastLinkNum,
		Data:        strg.Data,
	}

	// Преобразуем в JSON с отступами
	encoded, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}

	// Записываем в файл (файл будет создан, если его нет)
	return strg.reader.WriteFile(strg.filePath, encoded, 0644)
}

// GenerateID увеличивает счётчик и возвращает новый номер.
func (strg *Storage) GenerateID() int {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	strg.LastLinkNum++
	return strg.LastLinkNum
}

// AddRecord сохраняет новый результат по ID.
func (strg *Storage) AddRecord(id int, data map[string]string) {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	strg.Data[id] = data
}

// GetRecords возвращает данные по указанным ID групп ссылок.
func (strg *Storage) GetRecords(ids []int) map[int]map[string]string {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	out := make(map[int]map[string]string)

	for _, id := range ids {
		if val, ok := strg.Data[id]; ok {
			out[id] = val
		}
	}

	return out
}
