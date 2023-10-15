package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kupriyanovkk/shortener/internal/models"
	_ "github.com/lib/pq"
)

type StorageModel interface {
	GetValue(ctx context.Context, key string) (string, error)
	AddValue(ctx context.Context, key string, value string)
	Ping() error
	Bootstrap(ctx context.Context) error
}

const (
	DataBaseStorage   = 0
	FileWriterStorage = 1
	MemoryStorage     = 2
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

type Storage struct {
	mode   int
	values map[string]string
	file   *os.File
	writer *bufio.Writer
	db     *sql.DB
}

func (s Storage) Ping() error {
	err := s.db.Ping()
	return err
}

func (s Storage) Bootstrap(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS shortener(
			uuid serial PRIMARY KEY,
			short varchar(128),
			original TEXT
		)
	`)

	return tx.Commit()
}

func (s Storage) FindOriginalURL(ctx context.Context, short string) (originalURL string, err error) {
	var original sql.NullString
	row := s.db.QueryRowContext(ctx, `SELECT original FROM shortener WHERE short = $1`, short)
	err = row.Scan(&original)
	return original.String, err
}

func (s Storage) GetValue(ctx context.Context, key string) (string, error) {
	if s.mode == DataBaseStorage {
		value, err := s.FindOriginalURL(ctx, key)
		return value, err
	}

	if value, ok := s.values[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", key)
}

func (s Storage) AddValue(ctx context.Context, key string, value string) {
	s.values[key] = value

	if s.mode == DataBaseStorage {
		err := s.SaveURL(ctx, key, value)
		if err != nil {
			log.Fatal(`SaveURL `, err)
		}
	}

	if s.mode == FileWriterStorage {
		uuid += 1
		if err := s.WriteValue(&models.URL{
			UUID:     uuid,
			Short:    key,
			Original: value,
		}); err != nil {
			log.Fatal(`WriteValue `, err)
		}
	}
}

func (s Storage) SaveURL(ctx context.Context, key string, value string) error {
	_, err := s.db.ExecContext(ctx, `
			INSERT INTO shortener
			(short, original)
			VALUES
			($1, $2);
	`, key, value)

	return err
}

func (s *Storage) WriteValue(value *models.URL) error {
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

func getDBStorage(dbDSN string) Storage {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}

	newStorage := Storage{
		mode:   DataBaseStorage,
		db:     db,
		values: make(map[string]string),
	}
	newStorage.Bootstrap(context.Background())

	return newStorage
}

func getFileStorage(filename string) Storage {
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

	return Storage{
		mode:   FileWriterStorage,
		values: values,
		file:   file,
		writer: bufio.NewWriter(file),
	}
}

func NewStorage(filename, dbDSN string) Storage {
	if dbDSN != "" {
		return getDBStorage(dbDSN)
	}

	if filename != "" {
		return getFileStorage(filename)
	}

	return Storage{
		mode:   MemoryStorage,
		values: make(map[string]string),
	}
}
