package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/invine/portfolio/internal/domain/portfolio"
	_ "github.com/mattn/go-sqlite3"
)

type SQLitePortfolioRepository struct {
	db *sql.DB
}

type querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type transactionModel struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PortfolioID uuid.UUID
	Asset       string
	Quantity    int
	Price       float64
	DateString  string
}

type portfolioModel struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

func NewSQLitePortfolioRepository(db *sql.DB) (*SQLitePortfolioRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database required")
	}

	sqlStmt := `
        CREATE TABLE IF NOT EXISTS portfolios
        (
            id text not null primary key,
            userid text not null,
            name text
        );
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("can't insert table: %w", err)
	}

	sqlStmt = `
        CREATE TABLE IF NOT EXISTS transactions
        (
            id text not null primary key,
            userid text not null,
            portfolioid text not null,
            date text not null,
            asset text not null,
            price real not null,
            quantity integer not null
        );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("can't insert table: %w", err)
	}

	r := &SQLitePortfolioRepository{db: db}
	return r, nil
}

func (r *SQLitePortfolioRepository) CreatePortfolio(ctx context.Context, p *portfolio.Portfolio) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't create portfolio: %w", err)
	}
	defer tx.Rollback()

	pm := portfolioToPortfolioModel(p)

	sqlStmt := `insert into portfolios (id, userid, name) values ($1, $2, $3)`
	if _, err := tx.ExecContext(ctx, sqlStmt, pm.ID, pm.UserID, pm.Name); err != nil {
		return fmt.Errorf("can't create portfolio: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't create portfolio: %w", err)
	}

	return nil
}

func (r *SQLitePortfolioRepository) GetPortfolio(ctx context.Context, userID, id uuid.UUID) (*portfolio.Portfolio, error) {
	// TODO combine with UpdatePortfolio and nmove to function
	pm, err := r.getPortfolio(ctx, r.db, userID, id, false)
	if err != nil {
		return nil, fmt.Errorf("can't find portfolio with id %s: %w", id.String(), err)
	}

	trms, err := r.getAllTransactions(ctx, r.db, pm.UserID, pm.ID, false)
	if err != nil {
		return nil, fmt.Errorf("can't find portfolio with id %s: %w", id.String(), err)
	}
	trs, err := transactionModelToTransactions(trms)
	if err != nil {
		return nil, fmt.Errorf("can't find portfolio with id %s: %w", id.String(), err)
	}

	p, err := portfolio.NewPortfolio(pm.ID, pm.UserID, pm.Name, trs)
	if err != nil {
		return nil, fmt.Errorf("can't find portfolio with id %s: %w", id.String(), err)
	}

	return p, nil
}

func (r *SQLitePortfolioRepository) GetAllPortfolios(ctx context.Context, userID uuid.UUID) ([]*portfolio.Portfolio, error) {
	pms, err := r.getAllPortfolios(ctx, r.db, userID, false)
	if err != nil {
		return nil, fmt.Errorf("can't list portfolios for user %s: %w", userID.String(), err)
	}

	portfolios := []*portfolio.Portfolio{}
	for _, pm := range pms {
		p, err := portfolio.NewPortfolio(pm.ID, pm.UserID, pm.Name, nil)
		if err != nil {
			return nil, fmt.Errorf("can't list portfolio %s for user %s: %w", pm.ID, userID.String(), err)
		}
		portfolios = append(portfolios, p)
	}

	return portfolios, nil
}

func (r *SQLitePortfolioRepository) GetAllTransactions(ctx context.Context, userID, portfolioID uuid.UUID) ([]*portfolio.Transaction, error) {
	trms, err := r.getAllTransactions(ctx, r.db, userID, portfolioID, false)
	if err != nil {
		return nil, fmt.Errorf("can't list transactions for portfolio %s: %w", portfolioID.String(), err)
	}

	transactions, err := transactionModelToTransactions(trms)
	if err != nil {
		return nil, fmt.Errorf("can't list transactions for portfolio %s: %w", portfolioID.String(), err)
	}

	return transactions, nil
}

func (r *SQLitePortfolioRepository) UpdatePortfolio(ctx context.Context, userID, id uuid.UUID, updateFn func(p *portfolio.Portfolio) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}
	defer tx.Rollback()

	// TODO combine with GetPortfolio and nmove to function
	pm, err := r.getPortfolio(ctx, tx, userID, id, true)
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}
	trms, err := r.getAllTransactions(ctx, r.db, pm.UserID, pm.ID, false)
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}
	trs, err := transactionModelToTransactions(trms)
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}
	p, err := portfolio.NewPortfolio(pm.ID, pm.UserID, pm.Name, trs)
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	if err := updateFn(p); err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	pm = portfolioToPortfolioModel(p)
	trms = portfolioToTransactionModel(p)

	sqlStmt := "update portfolios set name = $3 where userid = $1 and id = $2"
	if _, err := tx.ExecContext(ctx, sqlStmt, pm.UserID, pm.ID, pm.Name); err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	err = r.upsertTransactions(ctx, tx, trms)
	if err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	return nil
}

func (r *SQLitePortfolioRepository) DeletePortfolio(ctx context.Context, userID, id uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't delete portfolio %s: %w", id.String(), err)
	}
	defer tx.Rollback()

	sqlStmt := "DELETE FROM transactions WHERE userid=$1 AND portfolioid=$2"
	if _, err := tx.ExecContext(ctx, sqlStmt, userID, id); err != nil {
		return fmt.Errorf("can't delete portfolio %s: %w", id.String(), err)
	}

	sqlStmt = "DELETE FROM portfolios WHERE userid=$1 AND id=$2"
	if _, err := tx.ExecContext(ctx, sqlStmt, userID, id); err != nil {
		return fmt.Errorf("can't delete portfolio %s: %w", id.String(), err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't update portfolio %s: %w", id.String(), err)
	}

	return nil
}

func (r *SQLitePortfolioRepository) upsertTransactions(ctx context.Context, tx *sql.Tx, trms []*transactionModel) error {
	sqlStmt := `
        INSERT INTO
        transactions(id, userid, portfolioid, date, asset, price, quantity)
        VALUES($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT(id) DO UPDATE SET
        date=excluded.date, asset=excluded.asset, price=excluded.price, quantity=excluded.quantity
	`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("can't upsert transactions: %w", err)
	}
	defer stmt.Close()

	for _, trm := range trms {
		if _, err := stmt.ExecContext(ctx, trm.ID, trm.UserID, trm.PortfolioID, trm.DateString, trm.Asset, trm.Price, trm.Quantity); err != nil {
			return fmt.Errorf("can't upsert transaction %s: %w", trm.ID.String(), err)
		}
	}

	return nil
}

func (r *SQLitePortfolioRepository) getPortfolio(ctx context.Context, db rowQuerier, userID, id uuid.UUID, forUpdate bool) (*portfolioModel, error) {
	sqlStmt := "select name from portfolios where userid = $1 and id = $2"
	// if forUpdate {
	// 	sqlStmt += " for update"
	// }
	row := db.QueryRowContext(ctx, sqlStmt, userID, id)

	var name string
	if err := row.Scan(&name); err != nil {
		return nil, fmt.Errorf("portfolio %s not found: %w", id.String(), err)
	}

	pm := &portfolioModel{
		ID:     id,
		UserID: userID,
		Name:   name,
	}

	return pm, nil
}

func (r *SQLitePortfolioRepository) getAllPortfolios(ctx context.Context, db querier, userID uuid.UUID, forUpdate bool) ([]*portfolioModel, error) {
	sqlStmt := `select id, name from portfolios where userid = $1`
	// if forUpdate {
	// 	sqlStmt += " for update"
	// }
	rows, err := db.QueryContext(ctx, sqlStmt, userID)
	if err != nil {
		return nil, fmt.Errorf("can't list portfolios for user %s: %w", userID.String(), err)
	}
	defer rows.Close()

	pms := []*portfolioModel{}
	for rows.Next() {
		var (
			id   uuid.UUID
			name string
		)
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("can't list portfolios for user %s: %w", userID.String(), err)
		}
		pms = append(pms, &portfolioModel{
			ID:     id,
			UserID: userID,
			Name:   name,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("can't list portfolios for user %s: %w", userID.String(), err)
	}

	return pms, nil
}

func (r *SQLitePortfolioRepository) getAllTransactions(ctx context.Context, db querier, userID, portfolioID uuid.UUID, forUpdate bool) ([]*transactionModel, error) {
	sqlStmt := `select id, asset, quantity, price, date from transactions where userid = $1 and portfolioid = $2`
	// if forUpdate {
	// 	sqlStmt += " for update"
	// }
	rows, err := db.QueryContext(ctx, sqlStmt, userID, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("can't list transactions for portfolio %s: %w", portfolioID.String(), err)
	}
	defer rows.Close()

	trms := []*transactionModel{}
	for rows.Next() {
		var (
			id         uuid.UUID
			asset      string
			quantity   int
			price      float64
			dateString string
		)
		err := rows.Scan(
			&id,
			&asset,
			&quantity,
			&price,
			&dateString,
		)
		if err != nil {
			return nil, fmt.Errorf("can't list transactions for portfolio %s: %w", portfolioID.String(), err)
		}
		trms = append(trms, &transactionModel{
			ID:          id,
			UserID:      userID,
			PortfolioID: portfolioID,
			Asset:       asset,
			Quantity:    quantity,
			Price:       price,
			DateString:  dateString,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("can't list transactions for portfolio %s: %w", portfolioID.String(), err)
	}

	return trms, nil
}

func portfolioToPortfolioModel(p *portfolio.Portfolio) *portfolioModel {
	return &portfolioModel{
		ID:     p.ID(),
		UserID: p.UserID(),
		Name:   p.Name(),
	}
}

func portfolioToTransactionModel(p *portfolio.Portfolio) []*transactionModel {
	trms := []*transactionModel{}
	for _, t := range p.Transactions() {
		trms = append(trms, &transactionModel{
			ID:          t.ID(),
			UserID:      p.UserID(),
			PortfolioID: p.ID(),
			Asset:       t.Asset(),
			Quantity:    t.Quantity(),
			Price:       t.Price(),
			DateString:  t.Date().Format(time.RFC3339),
		})
	}
	return trms
}

func transactionModelToTransactions(trms []*transactionModel) ([]*portfolio.Transaction, error) {
	trs := []*portfolio.Transaction{}
	for _, trm := range trms {
		date, err := time.Parse(time.RFC3339, trm.DateString)
		if err != nil {
			return nil, fmt.Errorf("incorrect transaction parameter: %w", err)
		}
		tr, err := portfolio.NewTransaction(trm.ID, date, trm.Asset, trm.Quantity, trm.Price)
		if err != nil {
			return nil, fmt.Errorf("incorrect transaction parameter: %w", err)
		}
		trs = append(trs, tr)
	}
	return trs, nil
}
