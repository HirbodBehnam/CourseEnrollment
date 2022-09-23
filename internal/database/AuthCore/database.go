package AuthCore

import (
	"CourseEnrollment/pkg/course"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	db *pgxpool.Pool
}

// AuthUser will authorize the user
func (db Database) AuthUser(ctx context.Context, id uint64, password string, isStaff bool) (bool, course.DepartmentID, error) {
	// Query data
	var hashedPassword string
	var deparmentID course.DepartmentID
	var err error
	if isStaff {
		err = db.db.QueryRow(ctx, "SELECT password, department_id FROM staff WHERE id=$1", id).Scan(&hashedPassword, &deparmentID)
	} else {
		err = db.db.QueryRow(ctx, "SELECT password, department_id FROM students WHERE id=$1", id).Scan(&hashedPassword, &deparmentID)
	}
	if err != nil {
		return false, 0, errors.Wrap(err, "cannot query row")
	}
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, deparmentID, nil
}
