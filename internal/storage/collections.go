package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kirre02/net-cli/internal/domain"
)

var errorCollectionNotFound = errors.New("collection not found")

func (s *SQLite) CreateCollection(ctx context.Context, collection domain.Collection) error {
	if err := collection.Validate(); err != nil {
		return err
	}

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO collections (id, name) VALUES (?, ?)`,
		collection.ID,
		collection.Name,
	)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	return nil
}

func (s *SQLite) GetCollection(ctx context.Context, id string) (domain.Collection, error) {
	var collection domain.Collection

	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, name FROM collections WHERE id = ?`,
		id,
	).Scan(&collection.ID, &collection.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Collection{}, errorCollectionNotFound
	}
	if err != nil {
		return domain.Collection{}, fmt.Errorf("get collection: %w", err)
	}

	return collection, nil
}

func (s *SQLite) GetCollections(ctx context.Context) ([]domain.Collection, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name FROM collections ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("get collections: %w", err)
	}
	defer rows.Close()

	var collections []domain.Collection

	for rows.Next() {
		var collection domain.Collection

		if err := rows.Scan(&collection.ID, &collection.Name); err != nil {
			return nil, fmt.Errorf("scan collection: %w", err)
		}

		collections = append(collections, collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate collections: %w", err)
	}

	return collections, nil
}

func (s *SQLite) UpdateCollection(ctx context.Context, collection domain.Collection) error {
	if err := collection.Validate(); err != nil {
		return err
	}

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE collections SET name = ? WHERE id = ?`,
		collection.Name,
		collection.ID,
	)
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated collection: %w", err)
	}
	if affected == 0 {
		return errorCollectionNotFound
	}

	return nil
}

func (s *SQLite) DeleteCollection(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(
		ctx,
		`DELETE FROM collections WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted collection: %w", err)
	}
	if affected == 0 {
		return errorCollectionNotFound
	}

	return nil
}
