package AuthCore

import (
	"CourseEnrollment/pkg/course"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(db *pgxpool.Pool) Database {
	return Database{db}
}

// AuthUser will authorize the user
func (db Database) AuthUser(ctx context.Context, id uint64, password string, isStaff bool) (bool, course.DepartmentID, error) {
	// Query data
	var hashedPassword string
	var departmentID course.DepartmentID
	var err error
	if isStaff {
		err = db.db.QueryRow(ctx, "SELECT password, department_id FROM staff WHERE id=$1", id).Scan(&hashedPassword, &departmentID)
	} else {
		err = db.db.QueryRow(ctx, "SELECT password, department_id FROM students WHERE id=$1", id).Scan(&hashedPassword, &departmentID)
	}
	// No user found
	if err == pgx.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, errors.Wrap(err, "cannot query row")
	}
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, departmentID, nil
}

// Close will close the database connection
func (db Database) Close() {
	db.db.Close()
}
