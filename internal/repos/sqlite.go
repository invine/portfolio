package repos

import (
	"database/sql"
	"fmt"

	"github.com/invine/Portfolio/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type PortfolioRepo struct {
	db *sql.DB
}

func NewPortfolioRepo(DBName string) (*PortfolioRepo, error) {
	// os.Remove(DBName)

	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	sqlStmt := `
	create table assets (name text not null primary key, amount integer not null);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil && err.Error() != "table assets already exists" {
		return nil, fmt.Errorf("can't insert table: %w", err)
	}

	return &PortfolioRepo{db: db}, nil
}

func (p *PortfolioRepo) Close() error {
	return p.db.Close()

}

func (pr *PortfolioRepo) Read() (*models.Portfolio, error) {
	rows, err := pr.db.Query("select name, amount from assets")
	if err != nil {
		return new(models.Portfolio), fmt.Errorf("can't load portfolio: %w", err)
	}
	defer rows.Close()

	assets := map[string]int{}

	for rows.Next() {
		var name string
		var amount int
		err = rows.Scan(&name, &amount)
		if err != nil {
			return new(models.Portfolio), fmt.Errorf("can't load portfolio: %w", err)
		}
		assets[name] = amount
	}

	err = rows.Err()
	if err != nil {
		return new(models.Portfolio), fmt.Errorf("can't load portfolio: %w", err)
	}
	return models.NewPortfolio(assets), nil
}

func (pr *PortfolioRepo) Update(t models.Transaction) error {
	p, err := pr.Read()
	if err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}

	if err := p.Apply(t); err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}

	tx, err := pr.db.Begin()
	if err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}
	upsertStmt, err := tx.Prepare("insert into assets(name, amount) values(?, ?) on conflict(name) do update set amount=?")
	if err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}
	deleteStmt, err := tx.Prepare("delete from assets where name = ?")
	if err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}
	defer upsertStmt.Close()

	for k, v := range p.Assets() {
		if v == 0 {
			_, err = deleteStmt.Exec(k)
			if err != nil {
				return fmt.Errorf("can't update portfolio: %w", err)
			}
			continue
		}
		_, err = upsertStmt.Exec(k, v, v)
		if err != nil {
			return fmt.Errorf("can't update portfolio: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't update portfolio: %w", err)
	}

	return nil
}
