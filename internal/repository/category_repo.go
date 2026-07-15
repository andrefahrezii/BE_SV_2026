package repository

import (
	"database/sql"
	"time"

	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(cat *model.Category) error {
	query := `INSERT INTO sv_portal.categories (name, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING id`
	now := time.Now()
	return r.db.QueryRow(query, cat.Name, now, now).Scan(&cat.ID)
}

func (r *CategoryRepo) FindByID(id int) (*model.Category, error) {
	cat := &model.Category{}
	err := r.db.QueryRow(
		"SELECT id, name, created_at, updated_at FROM sv_portal.categories WHERE id=$1", id,
	).Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cat, err
}

func (r *CategoryRepo) FindByName(name string) (*model.Category, error) {
	cat := &model.Category{}
	err := r.db.QueryRow(
		"SELECT id, name, created_at, updated_at FROM sv_portal.categories WHERE name=$1", name,
	).Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cat, err
}

func (r *CategoryRepo) Update(id int, name string) error {
	_, err := r.db.Exec(
		"UPDATE sv_portal.categories SET name=$1, updated_at=NOW() WHERE id=$2", name, id,
	)
	return err
}

func (r *CategoryRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM sv_portal.categories WHERE id=$1", id)
	return err
}

func (r *CategoryRepo) List() ([]model.Category, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM sv_portal.categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}
