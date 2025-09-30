package repository

import (
	"cruder/internal/model"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	return db, mock

}

func TestGetAll_Success(t *testing.T) {
	// Given: a mock db with two users returned from query
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{"id", "username", "email", "full_name"}).
		AddRow(1, "john_doe", "john@doe.ee", "John Doe").
		AddRow(2, "jane_doe", "jane@doe.ee", "Jane Doe")
	mock.ExpectQuery(`SELECT id, username, email, full_name FROM users`).
		WillReturnRows(rows)

	// When: calling GetAll
	users, err := repo.GetAll()

	// Then: two users should be returned without error
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "john_doe", users[0].Username)
	assert.Equal(t, "jane_doe", users[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUsername_Success(t *testing.T) {
	// Given: a user with username "john_doe" exists
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	row := sqlmock.NewRows([]string{"id", "username", "email", "full_name"}).
		AddRow(1, "john_doe", "john@doe.ee", "John Doe")
	mock.ExpectQuery(`SELECT id, username, email, full_name FROM users WHERE username = \$1`).
		WithArgs("john_doe").
		WillReturnRows(row)

	// When: calling GetByUsername with "john_doe"
	user, err := repo.GetByUsername("john_doe")

	// Then: the user should be returned with matching ID and username
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "john_doe", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUsername_NotFound(t *testing.T) {
	// Given: no user exists with username "missing"
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(`SELECT id, username, email, full_name FROM users WHERE username = \$1`).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	// When: calling GetByUsername with "missing"
	user, err := repo.GetByUsername("missing")

	// Then: ErrUserNotFound should be returned and user should be nil
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	// Given: a user with ID 1 exists
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	row := sqlmock.NewRows([]string{"id", "username", "email", "full_name"}).
		AddRow(1, "john_doe", "john@doe.ee", "John Doe")
	mock.ExpectQuery(`SELECT id, username, email, full_name FROM users WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(row)

	// When: calling GetByID with 1
	user, err := repo.GetByID(1)

	// Then: the user should be returned without error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Given: no user exists with ID 99
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(`SELECT id, username, email, full_name FROM users WHERE id = \$1`).
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	// When: calling GetByID with 99
	user, err := repo.GetByID(99)

	// Then: ErrUserNotFound should be returned and user should be nil
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Success(t *testing.T) {
	// Given: a new user to be inserted successfully
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	newUser := &model.User{Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	row := sqlmock.NewRows([]string{"id", "username", "email", "full_name"}).
		AddRow(1, newUser.Username, newUser.Email, newUser.FullName)
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(newUser.Username, newUser.Email, newUser.FullName).
		WillReturnRows(row)

	// When: calling Create with newUser
	created, err := repo.Create(newUser)

	// Then: user should be returned with ID set
	assert.NoError(t, err)
	assert.Equal(t, int64(1), created.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	// Given: inserting a user with duplicate username
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	newUser := &model.User{Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	pqErr := &pq.Error{Code: "23505", Constraint: "users_username_key"}
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(newUser.Username, newUser.Email, newUser.FullName).
		WillReturnError(pqErr)

	// When: calling Create with duplicate username
	created, err := repo.Create(newUser)

	// Then: ErrUsernameAlreadyExists should be returned
	assert.ErrorIs(t, err, ErrUsernameAlreadyExists)
	assert.Nil(t, created)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	// Given: inserting a user with duplicate email
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	newUser := &model.User{Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	pqErr := &pq.Error{Code: "23505", Constraint: "users_email_key"}
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(newUser.Username, newUser.Email, newUser.FullName).
		WillReturnError(pqErr)

	// When: calling Create with duplicate email
	created, err := repo.Create(newUser)

	// Then: ErrEmailAlreadyExists should be returned
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	assert.Nil(t, created)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DuplicateOther(t *testing.T) {
	// Given: inserting a user fails due to other duplicate constraint
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	newUser := &model.User{Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	pqErr := &pq.Error{Code: "23505", Constraint: "other_key"}
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(newUser.Username, newUser.Email, newUser.FullName).
		WillReturnError(pqErr)

	// When: calling Create with some other duplicate
	created, err := repo.Create(newUser)

	// Then: ErrUserAlreadyExists should be returned
	assert.ErrorIs(t, err, ErrUserAlreadyExists)
	assert.Nil(t, created)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_Success(t *testing.T) {
	// Given: a user with ID 1 exists and will be deleted
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`DELETE FROM users WHERE id = \$1 RETURNING id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// When: calling Delete with ID 1
	err := repo.Delete(1)

	// Then: no error should be returned
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_NotFound(t *testing.T) {
	// Given: no user exists with ID 99
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(`DELETE FROM users WHERE id = \$1 RETURNING id`).
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	// When: calling Delete with ID 99
	err := repo.Delete(99)

	// Then: ErrUserNotFound should be returned
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_Success(t *testing.T) {
	// Given: an existing user is updated successfully
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	user := &model.User{ID: 1, Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	row := sqlmock.NewRows([]string{"id", "username", "email", "full_name"}).
		AddRow(user.ID, user.Username, user.Email, user.FullName)
	mock.ExpectQuery(`UPDATE users SET username = \$1, email = \$2, full_name = \$3 WHERE id = \$4 RETURNING id, username, email, full_name`).
		WithArgs(user.Username, user.Email, user.FullName, user.ID).
		WillReturnRows(row)

	// When: calling Update with existing user
	updated, err := repo.Update(user)

	// Then: updated user should be returned without error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), updated.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_NotFound(t *testing.T) {
	// Given: no user exists with ID 99
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	user := &model.User{ID: 99, Username: "john_doe", Email: "john@doe.ee", FullName: "John Doe"}
	mock.ExpectQuery(`UPDATE users SET username = \$1, email = \$2, full_name = \$3 WHERE id = \$4 RETURNING id, username, email, full_name`).
		WithArgs(user.Username, user.Email, user.FullName, user.ID).
		WillReturnError(sql.ErrNoRows)

	// When: calling Update with non-existing user
	updated, err := repo.Update(user)

	// Then: ErrUserNotFound should be returned
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, updated)
	assert.NoError(t, mock.ExpectationsWereMet())
}
