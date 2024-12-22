// Flow: NewPostgresRepository -> postgresRepository instance created -> Subsequent methods called on the instance
package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()																			// Close the connection to the database
	PutAccount(ctx context.Context, a Account) error								// PutAccount stores the account in the database
	GetAccountByID(ctx context.Context, id string) (*Account, error)				// GetAccountByID retrieves an account by its ID
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)	// ListAccounts retrieves a list of accounts
}

type postgresRepository struct {								// postgresRepository is a PostgreSQL implementation of the Repository interface
	db *sql.DB
}
// above is not necessary but recommended and without it our function would be like this:
// func GetAccountByID(ctx context.Context, db *sql.DB, id string) (*Account, error) {

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)						// Open a new connection to the db
	if err != nil {
		return nil, err
	}
	err = db.Ping()												// Ping the db to check if the connection is successful
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)					// Select single row with the given id
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {													// Scan the filtered if the row exists
		return nil, err
	}
	return a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",							// Select all rows with pagination and take limit of rows
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()											// let rows run until the end of the function

	accounts := []Account{}
	for rows.Next() {											// Iterate over the rows until there are no more rows
		a := &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}