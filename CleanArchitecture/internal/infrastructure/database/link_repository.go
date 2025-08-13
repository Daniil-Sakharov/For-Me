package database

import (
	"context"
	"database/sql"
	"clean-url-shortener/internal/domain/link"
)

// LinkRepository реализует интерфейс link.Repository для PostgreSQL
type LinkRepository struct {
	db *DB
}

// NewLinkRepository создает новый репозиторий ссылок
func NewLinkRepository(db *DB) link.Repository {
	return &LinkRepository{
		db: db,
	}
}

// Save сохраняет ссылку в базе данных
func (r *LinkRepository) Save(ctx context.Context, l *link.Link) error {
	if l.ID == 0 {
		return r.create(ctx, l)
	}
	return r.update(ctx, l)
}

// create создает новую ссылку
func (r *LinkRepository) create(ctx context.Context, l *link.Link) error {
	query := `
		INSERT INTO links (original_url, short_code, user_id, clicks_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		l.OriginalURL,
		l.ShortCode,
		l.UserID,
		l.ClicksCount,
		l.CreatedAt,
		l.UpdatedAt,
	).Scan(&l.ID)

	return err
}

// update обновляет существующую ссылку
func (r *LinkRepository) update(ctx context.Context, l *link.Link) error {
	query := `
		UPDATE links 
		SET original_url = $1, short_code = $2, clicks_count = $3, updated_at = $4
		WHERE id = $5`

	_, err := r.db.ExecContext(
		ctx,
		query,
		l.OriginalURL,
		l.ShortCode,
		l.ClicksCount,
		l.UpdatedAt,
		l.ID,
	)

	return err
}

// FindByID находит ссылку по ID
func (r *LinkRepository) FindByID(ctx context.Context, id uint) (*link.Link, error) {
	query := `
		SELECT id, original_url, short_code, user_id, clicks_count, created_at, updated_at
		FROM links
		WHERE id = $1`

	l := &link.Link{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.OriginalURL,
		&l.ShortCode,
		&l.UserID,
		&l.ClicksCount,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return l, nil
}

// FindByShortCode находит ссылку по короткому коду
func (r *LinkRepository) FindByShortCode(ctx context.Context, shortCode string) (*link.Link, error) {
	query := `
		SELECT id, original_url, short_code, user_id, clicks_count, created_at, updated_at
		FROM links
		WHERE short_code = $1`

	l := &link.Link{}
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(
		&l.ID,
		&l.OriginalURL,
		&l.ShortCode,
		&l.UserID,
		&l.ClicksCount,
		&l.CreatedAt,
		&l.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return l, nil
}

// FindByUserID находит все ссылки пользователя с пагинацией
func (r *LinkRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*link.Link, error) {
	query := `
		SELECT id, original_url, short_code, user_id, clicks_count, created_at, updated_at
		FROM links
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*link.Link
	for rows.Next() {
		l := &link.Link{}
		err := rows.Scan(
			&l.ID,
			&l.OriginalURL,
			&l.ShortCode,
			&l.UserID,
			&l.ClicksCount,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	return links, nil
}

// Delete удаляет ссылку по ID
func (r *LinkRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ExistsByShortCode проверяет существование ссылки с указанным коротким кодом
func (r *LinkRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&exists)
	return exists, err
}

// CountByUserID возвращает общее количество ссылок пользователя
func (r *LinkRepository) CountByUserID(ctx context.Context, userID uint) (int, error) {
	query := `SELECT COUNT(*) FROM links WHERE user_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// IncrementClicks увеличивает счетчик переходов по ссылке
func (r *LinkRepository) IncrementClicks(ctx context.Context, linkID uint) error {
	query := `UPDATE links SET clicks_count = clicks_count + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, linkID)
	return err
}