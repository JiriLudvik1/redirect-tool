package analytics

import (
	"database/sql"
	"errors"
	"fmt"
)

func (manager *DbManager) AddClientAndRedirect(remoteAddress string, targetURL string) error {
	// Begin a transaction
	tx, err := manager.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Println(err)
				return
			}
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err = tx.Commit()
		}
	}()

	var clientId int64
	clientId, err = manager.getClientID(remoteAddress, tx)
	if err != nil {
		return err
	}

	if clientId == -1 {
		clientId, err = manager.insertClient(remoteAddress, tx)
	}

	err = manager.insertRedirect(targetURL, clientId, tx)
	if err != nil {
		return err
	}

	return nil
}

func (manager *DbManager) getClientID(remoteAddress string, transaction *sql.Tx) (int64, error) {
	var clientID int64
	query := `SELECT id FROM clients WHERE remote_address = ?`
	err := transaction.QueryRow(query, remoteAddress).Scan(&clientID)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, nil
	}
	return clientID, err
}

func (manager *DbManager) insertClient(remoteAddress string, transaction *sql.Tx) (int64, error) {
	query := `INSERT INTO clients (remote_address) VALUES (?)`
	result, err := transaction.Exec(query, remoteAddress)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (manager *DbManager) insertRedirect(targetURL string, clientID int64, transaction *sql.Tx) error {
	query := `INSERT INTO redirects (target_url, client_id) VALUES (?, ?)`
	_, err := transaction.Exec(query, targetURL, clientID)
	return err
}
