package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"local.dev/doc-analyzer/internal/pkg/storage/repository/postgres"
)

func TestSaveFile(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository with the mock database
	repo := postgres.NewFileRepo(db)

	// Test case: successful save
	t.Run("Successful save", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectExec("INSERT INTO files").
			WithArgs("file123", "test.txt", "hash123", "files/test.txt").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the method
		err := repo.SaveFile(
			context.Background(),
			"file123",
			"test.txt",
			"hash123",
			"files/test.txt",
		)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: database error
	t.Run("Database error", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectExec("INSERT INTO files").
			WithArgs("file123", "test.txt", "hash123", "files/test.txt").
			WillReturnError(errors.New("database error"))

		// Call the method
		err := repo.SaveFile(
			context.Background(),
			"file123",
			"test.txt",
			"hash123",
			"files/test.txt",
		)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save file metadata")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetFileByID(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository with the mock database
	repo := postgres.NewFileRepo(db)

	// Test case: successful get
	t.Run("Successful get", func(t *testing.T) {
		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"name", "location"}).
			AddRow("test.txt", "files/test.txt")

		mock.ExpectQuery("SELECT name, location FROM files").
			WithArgs("file123").
			WillReturnRows(rows)

		// Call the method
		name, location, err := repo.GetFileByID(
			context.Background(),
			"file123",
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "test.txt", name)
		assert.Equal(t, "files/test.txt", location)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: not found
	t.Run("Not found", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectQuery("SELECT name, location FROM files").
			WithArgs("file123").
			WillReturnError(sql.ErrNoRows)

		// Call the method
		_, _, err := repo.GetFileByID(
			context.Background(),
			"file123",
		)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: database error
	t.Run("Database error", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectQuery("SELECT name, location FROM files").
			WithArgs("file123").
			WillReturnError(errors.New("database error"))

		// Call the method
		_, _, err := repo.GetFileByID(
			context.Background(),
			"file123",
		)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file by id")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetFileByHash(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository with the mock database
	repo := postgres.NewFileRepo(db)

	// Test case: successful get
	t.Run("Successful get", func(t *testing.T) {
		// Set up mock expectations
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow("file123")

		mock.ExpectQuery("SELECT id FROM files").
			WithArgs("hash123").
			WillReturnRows(rows)

		// Call the method
		id, err := repo.GetFileByHash(
			context.Background(),
			"hash123",
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "file123", id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: not found
	t.Run("Not found", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectQuery("SELECT id FROM files").
			WithArgs("hash123").
			WillReturnError(sql.ErrNoRows)

		// Call the method
		id, err := repo.GetFileByHash(
			context.Background(),
			"hash123",
		)

		// Assert
		assert.NoError(t, err) // No error, just no file with this hash
		assert.Equal(t, "", id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: database error
	t.Run("Database error", func(t *testing.T) {
		// Set up mock expectations
		mock.ExpectQuery("SELECT id FROM files").
			WithArgs("hash123").
			WillReturnError(errors.New("database error"))

		// Call the method
		_, err := repo.GetFileByHash(
			context.Background(),
			"hash123",
		)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file by hash")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
