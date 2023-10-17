package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgerrcode"
	"github.com/kupriyanovkk/shortener/internal/models"
	"github.com/lib/pq"
)

var ErrConflict = errors.New("data conflict")

type AddValueOptions struct {
	Original string
	BaseURL  string
	Short    string
}

type StorageModel interface {
	GetValue(ctx context.Context, short string) (string, error)
	AddValue(ctx context.Context, opts AddValueOptions) (string, error)
	Ping() error
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

type DatabaseConnection interface {
	Ping() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Storage struct {
	mode   int
	values map[string]string
	file   *os.File
	writer *bufio.Writer
	db     DatabaseConnection
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

	tx.ExecContext(ctx, "CREATE UNIQUE INDEX url_id ON shortener (original)")

	return tx.Commit()
}

func (s Storage) FindOriginalURL(ctx context.Context, short string) (originalURL string, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT original FROM shortener WHERE short = $1`, short)
	err = row.Scan(&originalURL)
	return
}

func (s Storage) FindShortURL(ctx context.Context, original string) (shortURL string, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT short FROM shortener WHERE original = $1`, original)
	err = row.Scan(&shortURL)
	return
}

func (s Storage) GetValue(ctx context.Context, short string) (string, error) {
	if s.mode == DataBaseStorage {
		value, err := s.FindOriginalURL(ctx, short)
		return value, err
	}

	if value, ok := s.values[short]; ok {
		return value, nil
	}

	return "", fmt.Errorf("value doesn't exist by key %s", short)
}

func (s Storage) AddValue(ctx context.Context, opts AddValueOptions) (string, error) {
	if opts.Original == "" {
		return "", errors.New("original URL cannot be empty")
	}

	s.values[opts.Short] = opts.Original

	switch s.mode {
	case DataBaseStorage:
		err := s.SaveURL(ctx, opts.Short, opts.Original)

		if err != nil && errors.Is(err, ErrConflict) {
			short, _ := s.FindShortURL(ctx, opts.Original)
			result := fmt.Sprintf("%s/%s", opts.BaseURL, short)

			return result, err
		}
	case FileWriterStorage:
		result := fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short)
		uuid += 1
		if err := s.WriteValue(&models.URL{
			UUID:     uuid,
			Short:    opts.Short,
			Original: opts.Original,
		}); err != nil {
			return result, err
		}
	}

	return fmt.Sprintf("%s/%s", opts.BaseURL, opts.Short), nil
}

func (s Storage) SaveURL(ctx context.Context, short, original string) error {
	_, err := s.db.ExecContext(ctx, `
			INSERT INTO shortener
			(short, original)
			VALUES
			($1, $2);
	`, short, original)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = ErrConflict
		}
	}

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
