package infile

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/kupriyanovkk/shortener/internal/store"
)

var uuid = 0

func ReadValuesFromFile(scanner *bufio.Scanner) (map[string]string, error) {
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	values := make(map[string]string)
	for scanner.Scan() {
		value := models.URL{}
		err := json.Unmarshal(scanner.Bytes(), &value)
		if err != nil {
			return nil, err
		}
		uuid = value.UUID
		values[value.Short] = value.Original
	}

	return values, nil
}

type Store struct {
	values map[string]string
	file   *os.File
	writer *bufio.Writer
}

func (s Store) GetValue(ctx context.Context, short string) (string, error) {
	if value, ok := s.values[short]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

func (s Store) AddValue(ctx context.Context, opts store.AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = opts.Original
	result := fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short)
	uuid += 1
	if err := s.WriteValue(&models.URL{
		UUID:     uuid,
		Short:    opts.Short,
		Original: opts.Original,
	}); err != nil {
		return result, err
	}

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

func (s *Store) WriteValue(value *models.URL) error {
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

func (s Store) Ping() error {
	return nil
}

func NewStore(filename string) store.Store {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	values, readErr := ReadValuesFromFile(scanner)
	if readErr != nil {
		panic(readErr)
	}

	if len(values) == 0 {
		values = make(map[string]string)
	}

	return Store{
		values: values,
		file:   file,
		writer: bufio.NewWriter(file),
	}
}
