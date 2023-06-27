package user

import (
	"github.com/jmoiron/sqlx"
)

// InsertUserPayload represents the data required to create a user.
type InsertUserPayload struct {
	FirstName    string
	LastName     string
	PasswordHash string
	Email        string
}

// UserRow defines the marshalled struct of a User row from the database
type UserRow struct {
	Id           int64  `db:"id"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	PasswordHash string `db:"password_hash"`
	Email        string `db:"email"`
}

// FindUserByEmailPayload defines the input param struct for the FindUserByEmail func
type FindUserByEmailPayload struct {
	Email string
}

// UserService provides a persistence layer to interact with the database
type UserService struct {
	Db *sqlx.DB
}

// InsertUser nserts a new user into the database. If succesful, returns the inserted record ID
func (service UserService) InsertUser(payload InsertUserPayload) (int64, error) {
	result, err := service.Db.Exec("INSERT INTO user (first_name, last_name, email, password_hash, row_inserted) VALUES (?, ?, ?, ?, NOW())", payload.FirstName, payload.LastName, payload.Email, payload.PasswordHash)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindUserByEmail searches for a specific user record by email
func (service UserService) FindUserByEmail(payload FindUserByEmailPayload) ([]UserRow, error) {
	marshalledRows := []UserRow{}

	rows, err := service.Db.Queryx("SELECT id, first_name, last_name, email, password_hash FROM user WHERE email = ?", payload.Email)
	if err != nil {
		return marshalledRows, err
	}

	for rows.Next() {
		var marshalled UserRow
		err := rows.StructScan(&marshalled)
		if err != nil {
			return marshalledRows, err
		}
		marshalledRows = append(marshalledRows, marshalled)
	}

	return marshalledRows, nil
}
