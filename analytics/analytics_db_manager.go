package analytics

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DbManager struct {
	db *sql.DB
}

func NewAnalyticsDbManager(dataSourceName string) (*DbManager, error) {
	_, err := os.Stat(dataSourceName)
	if os.IsNotExist(err) {
		f, err := os.Create(dataSourceName)
		if err != nil {
			return nil, err
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(f)
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DbManager{db: db}, nil
}

func (manager *DbManager) RunMigrations() error {
	clientTableQuery := `
    CREATE TABLE IF NOT EXISTS clients (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        remote_address TEXT NOT NULL
    );
    `
	_, err := manager.db.Exec(clientTableQuery)
	if err != nil {
		return err
	}

	redirectTableQuery := `
    CREATE TABLE IF NOT EXISTS redirects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        target_url TEXT NOT NULL,
        client_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
    );
    `
	_, err = manager.db.Exec(redirectTableQuery)
	if err != nil {
		return err
	}

	indexQuery := `
    CREATE INDEX IF NOT EXISTS idx_remote_address ON clients (remote_address);
    `
	_, err = manager.db.Exec(indexQuery)

	return err
}

func (manager *DbManager) Close() error {
	return manager.db.Close()
}
