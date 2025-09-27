package repository

import (
	"context"
	"cruder/internal/model"
	"database/sql"
	"errors"
	"net/http"
	"os/exec"

	"github.com/lib/pq"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Delete(id int64) error
	Update(user *model.User) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]model.User, error) {
	rows, err := r.db.QueryContext(context.Background(), `SELECT id, username, email, full_name FROM users`)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	if err := r.db.QueryRowContext(context.Background(), `SELECT id, username, email, full_name FROM users WHERE username = $1`, username).
		Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int64) (*model.User, error) {
	var u model.User
	if err := r.db.QueryRowContext(context.Background(), `SELECT id, username, email, full_name FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.FullName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	// Need to catch unique constraint violations, right now they will return as 500 errors
	if err := r.db.QueryRowContext(context.Background(), `INSERT INTO users (username, email, full_name) VALUES ($1, $2, $3) RETURNING id, username, email, full_name`, user.Username, user.Email, user.FullName).
		Scan(&user.ID, &user.Username, &user.Email, &user.FullName); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				return nil, ErrUsernameAlreadyExists
			case "users_email_key":
				return nil, ErrEmailAlreadyExists
			default:
				return nil, ErrUserAlreadyExists
			}
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id int64) error {
	var idCheck int64
	if err := r.db.QueryRowContext(context.Background(), `DELETE FROM users WHERE id = $1 RETURNING id`, id).
		Scan(&idCheck); err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		} else {
			return err
		}
	}
	if idCheck == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.QueryRowContext(context.Background(), `UPDATE users SET username = $1, email = $2, full_name = $3 WHERE id = $4 RETURNING id, username, email, full_name`, user.Username, user.Email, user.FullName, user.ID).
		Scan(&user.ID, &user.Username, &user.Email, &user.FullName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}
	return user, nil
}

func InsecureHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ls", r.URL.Query().Get("dir"))
	cmd.Run()
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrEmailAlreadyExists = errors.New("email already exists")
