package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kirre02/net-cli/internal/domain"
)

var errorRequestNotFound = errors.New("request not found")

func (s *SQLite) CreateRequest(ctx context.Context, request domain.Request) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers, err := json.Marshal(request.Headers)
	if err != nil {
		return fmt.Errorf("marshal request headers: %w", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO requests (
			id, collection_id, name, method, url, headers, body
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		request.ID,
		request.CollectionID,
		request.Name,
		request.Method,
		request.URL,
		string(headers),
		request.Body,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	return nil
}

func (s *SQLite) GetRequest(ctx context.Context, id string) (domain.Request, error) {
	var request domain.Request
	var headers string

	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, collection_id, name, method, url, headers, body
		FROM requests WHERE id = ?`,
		id,
	).Scan(
		&request.ID,
		&request.CollectionID,
		&request.Name,
		&request.Method,
		&request.URL,
		&headers,
		&request.Body,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Request{}, errorRequestNotFound
	}
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request: %w", err)
	}

	if err := json.Unmarshal([]byte(headers), &request.Headers); err != nil {
		return domain.Request{}, fmt.Errorf("unmarshal request headers: %w", err)
	}

	return request, nil
}

func (s *SQLite) GetRequestsByCollection(ctx context.Context, collectionID string) ([]domain.Request, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, collection_id, name, method, url, headers, body
		FROM requests WHERE collection_id = ?
		ORDER BY name`,
		collectionID,
	)
	if err != nil {
		return nil, fmt.Errorf("get requests: %w", err)
	}
	defer rows.Close()

	var requests []domain.Request

	for rows.Next() {
		var request domain.Request
		var headers string

		if err := rows.Scan(
			&request.ID,
			&request.CollectionID,
			&request.Name,
			&request.Method,
			&request.URL,
			&headers,
			&request.Body,
		); err != nil {
			return nil, fmt.Errorf("scan request: %w", err)
		}

		if err := json.Unmarshal([]byte(headers), &request.Headers); err != nil {
			return nil, fmt.Errorf("unmarshal request headers: %w", err)
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate requests: %w", err)
	}

	return requests, nil
}

func (s *SQLite) UpdateRequest(ctx context.Context, request domain.Request) error {
	if err := request.Validate(); err != nil {
		return err
	}

	headers, err := json.Marshal(request.Headers)
	if err != nil {
		return fmt.Errorf("marshal request headers: %w", err)
	}

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE requests SET
			collection_id = ?,
			name = ?,
			method = ?,
			url = ?,
			headers = ?,
			body = ?
		WHERE id = ?`,
		request.CollectionID,
		request.Name,
		request.Method,
		request.URL,
		string(headers),
		request.Body,
		request.ID,
	)
	if err != nil {
		return fmt.Errorf("update request: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated request: %w", err)
	}
	if affected == 0 {
		return errorRequestNotFound
	}

	return nil
}

func (s *SQLite) DeleteRequest(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(
		ctx,
		`DELETE FROM requests WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted request: %w", err)
	}
	if affected == 0 {
		return errorRequestNotFound
	}

	return nil
}
