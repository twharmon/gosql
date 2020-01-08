package gosql

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	user := &User{
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^insert into user \(role, email, active\) values \(\?, \?, \?\)$`).WithArgs(user.Role, user.Email, user.Active).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DBConn.Insert(user); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// func TestSave(t *testing.T) {
// 	user := &User{
// 		ID:     5,
// 		Role:   "admin",
// 		Email:  "test@example.com",
// 		Active: true,
// 	}

// 	mock.ExpectExec(`^update user set role = \?, email = \?, active = \? where id = \?$`).WithArgs(user.Role, user.Email, user.Active, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

// 	if _, err := DBConn.Save(user); err != nil {
// 		t.Errorf("error was not expected while inserting: %s", err)
// 		return
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

func TestInsertAllTypes(t *testing.T) {
	all := &AllTypes{
		Uint64:      5,
		Int64:       5,
		Uint:        5,
		Int:         5,
		String:      "gopher",
		Float32:     5.5,
		Float64:     5.5,
		Blob:        []byte("gopher"),
		Time:        time.Now(),
		Bool:        true,
		NullBool:    NullBool{Valid: true, Bool: true},
		NullInt:     NullInt{Valid: true, Int: 5},
		NullInt64:   NullInt64{Valid: true, Int64: 5},
		NullInt32:   NullInt32{Valid: true, Int32: 5},
		NullUint:    NullUint{Valid: true, Uint: 5},
		NullUint64:  NullUint64{Valid: true, Uint64: 5},
		NullUint32:  NullUint32{Valid: true, Uint32: 5},
		NullString:  NullString{Valid: true, String: "gopher"},
		NullTime:    NullTime{Valid: true, Time: time.Now()},
		NullFloat32: NullFloat32{Valid: true, Float32: 5.5},
		NullFloat64: NullFloat64{Valid: true, Float64: 5.5},
	}
	mock.ExpectExec(`^insert into all_types \(uint64, uint, int, int64, float32, float64, string, blob, bool, time, null_string, null_uint64, null_uint32, null_int64, null_int32, null_int, null_uint, null_float64, null_float32, null_time, null_bool\) values \(\?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?\)$`).
		WithArgs(all.Uint64, all.Uint, all.Int, all.Int64, all.Float32, all.Float64, all.String, all.Blob, all.Bool, all.Time, all.NullString, all.NullUint64, all.NullUint32, all.NullInt64, all.NullInt32, all.NullInt, all.NullUint, all.NullFloat64, all.NullFloat32, all.NullTime, all.NullBool).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DBConn.Insert(all); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertAllTypesNull(t *testing.T) {
	all := &AllTypes{}
	mock.ExpectExec(`^insert into all_types \(uint64, uint, int, int64, float32, float64, string, blob, bool, time, null_string, null_uint64, null_uint32, null_int64, null_int32, null_int, null_uint, null_float64, null_float32, null_time, null_bool\) values \(\?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?, \?\)$`).
		WithArgs(all.Uint64, all.Uint, all.Int, all.Int64, all.Float32, all.Float64, all.String, all.Blob, all.Bool, all.Time, all.NullString, all.NullUint64, all.NullUint32, all.NullInt64, all.NullInt32, all.NullInt, all.NullUint, all.NullFloat64, all.NullFloat32, all.NullTime, all.NullBool).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DBConn.Insert(all); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	user := &User{
		ID:     5,
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^update user set role = \?, email = \?, active = \? where id = \?$`).WithArgs(user.Role, user.Email, user.Active, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	q := DBConn.ManualUpdate("user")
	q.Set("role = ?", user.Role)
	q.Set("email = ?", user.Email)
	q.Set("active = ?", user.Active)
	q.Where("id = ?", user.ID)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	mock.ExpectExec(`^delete from user where id = \?$`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	q := DBConn.ManualDelete("user")
	q.Where("id = ?", 1)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCount(t *testing.T) {
	rows := sqlmock.NewRows([]string{"count(*)"})
	rows.AddRow(10)
	mock.ExpectQuery(`^select count\(\*\) from user where id > \?$`).WithArgs(10).WillReturnRows(rows)

	q := DBConn.Count("user", "*")
	q.Where("id > ?", 10)

	if _, err := q.Exec(); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertPrimaryLast(t *testing.T) {
	user := &UserPrimaryLast{
		Role:   "admin",
		Email:  "test@example.com",
		Active: true,
	}

	mock.ExpectExec(`^insert into user_primary_last \(role, email, active\) values \(\?, \?, \?\)$`).WithArgs(user.Role, user.Email, user.Active).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DBConn.Insert(user); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertErrors(t *testing.T) {
	_, err := DBConn.Insert(User{})
	assertErr(t, "should return error if struct and not pointer to struct", err)

	_, err2 := DBConn.Insert(User{})
	assertErr(t, "should return error if pointer to non struct", err2)

	type Post struct {
		ID    int64
		Title string
	}
	_, err3 := DBConn.Insert(User{})
	assertErr(t, "should return error if struct not registered", err3)
}

func TestExec(t *testing.T) {
	mock.ExpectExec(`^update user set email = \? where id = \?$`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := DBConn.Exec("update user set email = ? where id = ?", 1); err != nil {
		t.Errorf("error was not expected while executing: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
