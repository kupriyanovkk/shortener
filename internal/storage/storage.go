package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var uuid = 0

func ReadValues(scanner *bufio.Scanner) (map[string]string, error) {
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	values := make(map[string]string)
	for scanner.Scan() {
		value := Value{}
		err := json.Unmarshal(scanner.Bytes(), &value)
		if err != nil {
			return nil, err
		}
		uuid = value.UUID
		values[value.ShortURL] = value.OriginalURL
	}

	return values, nil
}

type Value struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type StorageModel interface {
	GetValue(key string) (string, error)
	AddValue(key string, value string)
}

type Storage struct {
	values     map[string]string
	isWritable bool
	file       *os.File
	writer     *bufio.Writer
}

func (s Storage) GetValue(key string) (string, error) {
	if value, ok := s.values[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", key)
}

func (s Storage) AddValue(key string, value string) {
	s.values[key] = value

	if s.isWritable {
		uuid += 1
		if err := s.WriteValue(&Value{
			UUID:        uuid,
			ShortURL:    key,
			OriginalURL: value,
		}); err != nil {
			log.Fatal(`WriteValue `, err)
		}
	}
}

func (s *Storage) WriteValue(value *Value) error {
	data, err := json.Marshal(&value)

	if err != nil {
		return err
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}

	return s.writer.Flush()
}

func NewStorage(filename string) Storage {
	if filename != "" {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		values, readErr := ReadValues(scanner)
		if readErr != nil {
			panic(readErr)
		}

		if len(values) == 0 {
			values = make(map[string]string)
		}

		return Storage{
			values:     values,
			isWritable: true,
			file:       file,
			writer:     bufio.NewWriter(file),
		}
	}

	return Storage{
		values:     make(map[string]string),
		isWritable: false,
	}
}
