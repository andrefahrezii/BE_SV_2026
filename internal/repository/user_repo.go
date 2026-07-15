package repository

import (
	"database/sql"
	"time"

	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User) error {
	query := `INSERT INTO sv_portal.users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, user.Email, user.PasswordHash, user.Role, time.Now(), time.Now()).Scan(&user.ID)
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow("SELECT id, email, password_hash, role, created_at, updated_at FROM sv_portal.users WHERE email=$1", email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) GetByID(id int) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow("SELECT id, email, role, created_at, updated_at FROM sv_portal.users WHERE id=$1", id).
		Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) CountAdmins() (int, error) {
	var c int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM sv_portal.users WHERE role='admin'").Scan(&c); err != nil {
		return 0, err
	}
	return c, nil
}
