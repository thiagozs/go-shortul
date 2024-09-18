package database

import (
	"fmt"
	"log/slog"

	"github.com/thiagozs/go-shorturl/infra/database/memory"
	"github.com/thiagozs/go-shorturl/infra/database/sqlite"
)

type Kind int

const (
	Memory Kind = iota
	SQLite
)

func (k Kind) String() string {
	return [...]string{"Memory", "SQLite"}[k]
}

type DatabaseRepo interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string) (string, bool)
	GetStats(shortURL string) (string, bool)
	UpdateURL(shortURL, newOriginalURL string) error
	UpdateStats(shortURL, ip, referrer, geoLocation string) error
	Flush() (map[string]string, error)
	Backup() ([]byte, error)
	Import(data []byte) error
}

type Database struct {
	Kind   Kind
	Engine DatabaseRepo
	logger *slog.Logger
}

func NewDatabase(kind Kind, logger *slog.Logger) (*Database, error) {
	switch kind {
	case Memory:
		return &Database{Kind: kind,
			Engine: memory.NewURLStore(), logger: logger}, nil
	case SQLite:
		eng, err := sqlite.NewURLStore("./shorturl.db", logger)
		if err != nil {
			return nil, err
		}
		return &Database{Kind: kind,
			Engine: eng, logger: logger}, nil
	default:
		return &Database{}, fmt.Errorf("unsupported database kind")
	}
}

func (d *Database) Save(shortURL, originalURL string) {
	d.Engine.Save(shortURL, originalURL)
}

func (d *Database) Get(shortURL string) (string, bool) {
	return d.Engine.Get(shortURL)
}

func (d *Database) GetStats(shortURL string) (string, bool) {
	return d.Engine.GetStats(shortURL)
}

func (d *Database) UpdateURL(shortURL, newOriginalURL string) error {
	return d.Engine.UpdateURL(shortURL, newOriginalURL)
}

func (d *Database) UpdateStats(shortURL, ip, referrer, geoLocation string) error {
	return d.Engine.UpdateStats(shortURL, ip, referrer, geoLocation)
}

func (d *Database) Flush() (map[string]string, error) {
	return d.Engine.Flush()
}

func (d *Database) Backup() ([]byte, error) {
	return d.Engine.Backup()
}

func (d *Database) Import(data []byte) error {
	return d.Engine.Import(data)
}
