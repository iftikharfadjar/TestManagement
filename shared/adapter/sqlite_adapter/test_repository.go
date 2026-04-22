package sql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"boilerplate/services/test/domain"
)

type testRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) domain.TestRepository {
	return &testRepository{db: db}
}

func (r *testRepository) Create(ctx context.Context, test *domain.Test) (*domain.Test, error) {
	if test.TestID == "" {
		test.TestID = uuid.New().String()
	}

	now := time.Now()
	test.CreatedAt = now
	test.UpdatedAt = now

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tests (test_id, test_name, description, steps, status, remarks, created_by, created_at, updated_by, updated_at, is_active)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		test.TestID, test.TestName, test.Description, test.Steps, test.Status, test.Remarks,
		test.CreatedBy, test.CreatedAt, test.UpdatedBy, test.UpdatedAt, test.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (r *testRepository) GetByID(ctx context.Context, id string) (*domain.Test, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT test_id, test_name, description, steps, status, remarks, created_by, created_at, updated_by, updated_at, is_active
		 FROM tests WHERE test_id = ?`, id)

	var t domain.Test
	err := row.Scan(
		&t.TestID, &t.TestName, &t.Description, &t.Steps, &t.Status, &t.Remarks,
		&t.CreatedBy, &t.CreatedAt, &t.UpdatedBy, &t.UpdatedAt, &t.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("test not found")
		}
		return nil, err
	}

	return &t, nil
}

func (r *testRepository) List(ctx context.Context) ([]domain.Test, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT test_id, test_name, description, steps, status, remarks, created_by, created_at, updated_by, updated_at, is_active
		 FROM tests ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []domain.Test
	for rows.Next() {
		var t domain.Test
		if err := rows.Scan(
			&t.TestID, &t.TestName, &t.Description, &t.Steps, &t.Status, &t.Remarks,
			&t.CreatedBy, &t.CreatedAt, &t.UpdatedBy, &t.UpdatedAt, &t.IsActive,
		); err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}

	return tests, rows.Err()
}

func (r *testRepository) Update(ctx context.Context, id string, test *domain.Test) (*domain.Test, error) {
	test.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx,
		`UPDATE tests SET test_name = ?, description = ?, steps = ?, status = ?, remarks = ?, updated_by = ?, updated_at = ?, is_active = ?
		 WHERE test_id = ?`,
		test.TestName, test.Description, test.Steps, test.Status, test.Remarks,
		test.UpdatedBy, test.UpdatedAt, test.IsActive, id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("test not found")
	}

	return r.GetByID(ctx, id)
}

func (r *testRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tests WHERE test_id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("test not found")
	}

	return nil
}
