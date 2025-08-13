package database

import (
	"context"
	"database/sql"
	"clean-url-shortener/internal/domain/user"
)

// UserRepository реализует интерфейс user.Repository для PostgreSQL
// ЭТО РЕАЛИЗАЦИЯ ИНТЕРФЕЙСА ИЗ DOMAIN СЛОЯ В INFRASTRUCTURE СЛОЕ
type UserRepository struct {
	db *DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *DB) user.Repository {
	return &UserRepository{
		db: db,
	}
}

// Save сохраняет пользователя в базе данных
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	if u.ID == 0 {
		// Создаем нового пользователя
		return r.create(ctx, u)
	}
	// Обновляем существующего пользователя
	return r.update(ctx, u)
}

// create создает нового пользователя
func (r *UserRepository) create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (email, hashed_password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		u.Email,
		u.HashedPassword,
		u.Name,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

// update обновляет существующего пользователя
func (r *UserRepository) update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users 
		SET email = $1, hashed_password = $2, name = $3, updated_at = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.Email,
		u.HashedPassword,
		u.Name,
		u.UpdatedAt,
		u.ID,
	)

	return err
}

// FindByID находит пользователя по ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	query := `
		SELECT id, email, hashed_password, name, created_at, updated_at
		FROM users
		WHERE id = $1`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.HashedPassword,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail находит пользователя по email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, hashed_password, name, created_at, updated_at
		FROM users
		WHERE email = $1`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.HashedPassword,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Delete удаляет пользователя по ID
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Exists проверяет существование пользователя по email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// ПРИНЦИПЫ РЕАЛИЗАЦИИ REPOSITORY:
// 1. Реализует интерфейс из domain слоя
// 2. Содержит только SQL код и маппинг данных
// 3. Не содержит бизнес-логики
// 4. Обрабатывает ошибки базы данных
// 5. Преобразует доменные модели в SQL и обратно