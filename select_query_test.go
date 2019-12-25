package gosql

import (
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestBuilderSelectOneAllFields(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	rows.AddRow(control.ID, control.Role, control.Email, control.Active)
	mock.ExpectQuery(`^select \* from user where id = \? limit 1$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DBConn.Select("*").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestNullInt(t *testing.T) {
	control := &AllTypes{
		NullInt: NullInt{
			Valid: true,
			Int:   5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_int"})
	rows.AddRow(control.NullInt.Int)
	mock.ExpectQuery(`^select null_int from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullInt: NullInt{
			Valid: true,
			Int:   5,
		},
	}
	if err := DBConn.Select("null_int").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullInt != test.NullInt {
		t.Error(fmt.Errorf("test did not match control: control.NullInt %v, test.NullInt %v", control.NullInt, test.NullInt))
	}
}

func TestNullIntNull(t *testing.T) {
	control := &AllTypes{
		NullInt: NullInt{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_int"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_int from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullInt: NullInt{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_int").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullInt != test.NullInt {
		t.Error(fmt.Errorf("test did not match control: control.NullInt %v, test.NullInt %v", control.NullInt, test.NullInt))
	}
}

func TestNullInt64(t *testing.T) {
	control := &AllTypes{
		NullInt64: NullInt64{
			Valid: true,
			Int64: 5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_int64"})
	rows.AddRow(control.NullInt64.Int64)
	mock.ExpectQuery(`^select null_int64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullInt64: NullInt64{
			Valid: true,
			Int64: 5,
		},
	}
	if err := DBConn.Select("null_int64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullInt64 != test.NullInt64 {
		t.Error(fmt.Errorf("test did not match control: control.NullInt64 %v, test.NullInt64 %v", control.NullInt64, test.NullInt64))
	}
}

func TestNullInt64Null(t *testing.T) {
	control := &AllTypes{
		NullInt64: NullInt64{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_int64"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_int64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullInt64: NullInt64{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_int64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullInt64 != test.NullInt64 {
		t.Error(fmt.Errorf("test did not match control: control.NullInt64 %v, test.NullInt64 %v", control.NullInt64, test.NullInt64))
	}
}

func TestNullUint64(t *testing.T) {
	control := &AllTypes{
		NullUint64: NullUint64{
			Valid:  true,
			Uint64: 5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_uint64"})
	rows.AddRow(control.NullUint64.Uint64)
	mock.ExpectQuery(`^select null_uint64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullUint64: NullUint64{
			Valid:  true,
			Uint64: 5,
		},
	}
	if err := DBConn.Select("null_uint64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullUint64 != test.NullUint64 {
		t.Error(fmt.Errorf("test did not match control: control.NullUint64 %v, test.NullUint64 %v", control.NullUint64, test.NullUint64))
	}
}

func TestNullUint64Null(t *testing.T) {
	control := &AllTypes{
		NullUint64: NullUint64{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_uint64"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_uint64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullUint64: NullUint64{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_uint64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullUint64 != test.NullUint64 {
		t.Error(fmt.Errorf("test did not match control: control.NullUint64 %v, test.NullUint64 %v", control.NullUint64, test.NullUint64))
	}
}

func TestNullUint(t *testing.T) {
	control := &AllTypes{
		NullUint: NullUint{
			Valid: true,
			Uint:  5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_uint"})
	rows.AddRow(control.NullUint.Uint)
	mock.ExpectQuery(`^select null_uint from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullUint: NullUint{
			Valid: true,
			Uint:  5,
		},
	}
	if err := DBConn.Select("null_uint").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullUint != test.NullUint {
		t.Error(fmt.Errorf("test did not match control: control.NullUint %v, test.NullUint %v", control.NullUint, test.NullUint))
	}
}

func TestNullUintNull(t *testing.T) {
	control := &AllTypes{
		NullUint: NullUint{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_uint"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_uint from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullUint: NullUint{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_uint").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullUint != test.NullUint {
		t.Error(fmt.Errorf("test did not match control: control.NullUint %v, test.NullUint %v", control.NullUint, test.NullUint))
	}
}

func TestNullBool(t *testing.T) {
	control := &AllTypes{
		NullBool: NullBool{
			Valid: true,
			Bool:  true,
		},
	}
	rows := sqlmock.NewRows([]string{"null_bool"})
	rows.AddRow(control.NullBool.Bool)
	mock.ExpectQuery(`^select null_bool from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullBool: NullBool{
			Valid: true,
			Bool:  true,
		},
	}
	if err := DBConn.Select("null_bool").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullBool != test.NullBool {
		t.Error(fmt.Errorf("test did not match control: control.NullBool %v, test.NullBool %v", control.NullBool, test.NullBool))
	}
}

func TestNullBoolNull(t *testing.T) {
	control := &AllTypes{
		NullBool: NullBool{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_bool"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_bool from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullBool: NullBool{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_bool").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullBool != test.NullBool {
		t.Error(fmt.Errorf("test did not match control: control.NullBool %v, test.NullBool %v", control.NullBool, test.NullBool))
	}
}

func TestNullTime(t *testing.T) {
	now := time.Now()
	control := &AllTypes{
		NullTime: NullTime{
			Valid: true,
			Time:  now,
		},
	}
	rows := sqlmock.NewRows([]string{"null_time"})
	rows.AddRow(control.NullTime.Time)
	mock.ExpectQuery(`^select null_time from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullTime: NullTime{
			Valid: true,
			Time:  now,
		},
	}
	if err := DBConn.Select("null_time").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullTime != test.NullTime {
		t.Error(fmt.Errorf("test did not match control: control.NullTime %v, test.NullTime %v", control.NullTime, test.NullTime))
	}
}

func TestNullTimeNull(t *testing.T) {
	control := &AllTypes{
		NullTime: NullTime{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_time"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_time from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullTime: NullTime{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_time").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullTime != test.NullTime {
		t.Error(fmt.Errorf("test did not match control: control.NullTime %v, test.NullTime %v", control.NullTime, test.NullTime))
	}
}

func TestNullFloat32(t *testing.T) {
	control := &AllTypes{
		NullFloat32: NullFloat32{
			Valid:   true,
			Float32: 5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_float32"})
	rows.AddRow(control.NullTime.Time)
	mock.ExpectQuery(`^select null_float32 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullFloat32: NullFloat32{
			Valid:   true,
			Float32: 5,
		},
	}
	if err := DBConn.Select("null_float32").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullFloat32 != test.NullFloat32 {
		t.Error(fmt.Errorf("test did not match control: control.NullFloat32 %v, test.NullFloat32 %v", control.NullFloat32, test.NullFloat32))
	}
}

func TestNullFloat32Null(t *testing.T) {
	control := &AllTypes{
		NullFloat32: NullFloat32{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_float32"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_float32 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullFloat32: NullFloat32{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_float32").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullFloat32 != test.NullFloat32 {
		t.Error(fmt.Errorf("test did not match control: control.NullFloat32 %v, test.NullFloat32 %v", control.NullFloat32, test.NullFloat32))
	}
}

func TestNullFloat64(t *testing.T) {
	control := &AllTypes{
		NullFloat64: NullFloat64{
			Valid:   true,
			Float64: 5,
		},
	}
	rows := sqlmock.NewRows([]string{"null_float64"})
	rows.AddRow(control.NullTime.Time)
	mock.ExpectQuery(`^select null_float64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullFloat64: NullFloat64{
			Valid:   true,
			Float64: 5,
		},
	}
	if err := DBConn.Select("null_float64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullFloat64 != test.NullFloat64 {
		t.Error(fmt.Errorf("test did not match control: control.NullFloat64 %v, test.NullFloat64 %v", control.NullFloat64, test.NullFloat64))
	}
}

func TestNullFloat64Null(t *testing.T) {
	control := &AllTypes{
		NullFloat64: NullFloat64{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_float64"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_float64 from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullFloat64: NullFloat64{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_float64").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullFloat64 != test.NullFloat64 {
		t.Error(fmt.Errorf("test did not match control: control.NullFloat64 %v, test.NullFloat64 %v", control.NullFloat64, test.NullFloat64))
	}
}

func TestNullString(t *testing.T) {
	control := &AllTypes{
		NullString: NullString{
			Valid:  true,
			String: "gopher",
		},
	}
	rows := sqlmock.NewRows([]string{"null_string"})
	rows.AddRow(control.NullString.String)
	mock.ExpectQuery(`^select null_string from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullString: NullString{
			Valid:  true,
			String: "gopher",
		},
	}
	if err := DBConn.Select("null_string").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullString != test.NullString {
		t.Error(fmt.Errorf("test did not match control: control.NullString %v, test.NullString %v", control.NullString, test.NullString))
	}
}

func TestNullStringNull(t *testing.T) {
	control := &AllTypes{
		NullString: NullString{
			Valid: false,
		},
	}
	rows := sqlmock.NewRows([]string{"null_string"})
	rows.AddRow(nil)
	mock.ExpectQuery(`^select null_string from all_types limit 1$`).WillReturnRows(rows)
	test := &AllTypes{
		NullString: NullString{
			Valid: false,
		},
	}
	if err := DBConn.Select("null_string").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if control.NullString != test.NullString {
		t.Error(fmt.Errorf("test did not match control: control.NullString %v, test.NullString %v", control.NullString, test.NullString))
	}
}

func TestBuilderSelectOneSomeFields(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"email", "active"})
	rows.AddRow(control.Email, control.Active)
	mock.ExpectQuery(`^select email, active from user where id = \? limit 1$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DBConn.Select("email", "active").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	test.ID = control.ID
	test.Role = control.Role
	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneOneField(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"email"})
	rows.AddRow(control.Email)
	mock.ExpectQuery(`^select email from user where id = \? limit 1$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DBConn.Select("email").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	test.ID = control.ID
	test.Active = control.Active
	test.Role = control.Role
	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneJoin(t *testing.T) {
	control := makeUser()
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	rows.AddRow(control.ID, control.Role, control.Email, control.Active)
	mock.ExpectQuery(`^select user\.\* from user join post on post.user_id = user.id where user.id = \? limit 1$`).WithArgs(1).WillReturnRows(rows)

	test := new(User)
	if err := DBConn.Select("user.*").Where("user.id = ?", 1).Join("post on post.user_id = user.id").To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneExpanded(t *testing.T) {
	control := makeExpandedUser()
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	rows.AddRow(control.ID, control.Role, control.Email, control.Active)
	mock.ExpectQuery(`^select \* from expanded_user where id = \? limit 1$`).WithArgs(1).WillReturnRows(rows)

	test := new(ExpandedUser)
	if err := DBConn.Select("*").Where("id = ?", 1).To(test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertExpandedSame(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyAllFields(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for _, c := range control {
		rows.AddRow(c.ID, c.Role, c.Email, c.Active)
	}
	mock.ExpectQuery(`^select \* from user limit 1000$`).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("*").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}
func TestBuilderSelectManyValuesSomeFields(t *testing.T) {
	control := makeUserValuesSlice(3)
	rows := sqlmock.NewRows([]string{"email", "active"})
	for i := range control {
		rows.AddRow(control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select email, active from user limit 1000$`).WillReturnRows(rows)

	var test []User
	if err := DBConn.Select("email", "active").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
	}
	if err := assertSameSliceValues(control, test); err != nil {
		t.Error(err)
	}
}
func TestBuilderSelectManyOneField(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"email"})
	for i := range control {
		rows.AddRow(control[i].Email)
	}
	mock.ExpectQuery(`^select email from user limit 1000$`).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("email").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
		test[i].Active = control[i].Active
	}
	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyLimit(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for i := range control {
		rows.AddRow(control[i].ID, control[i].Role, control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select \* from user where id = \? limit 5$`).WithArgs(1).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("*").Where("id = ?", 1).Limit(5).To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyOrderBy(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for i := range control {
		rows.AddRow(control[i].ID, control[i].Role, control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select \* from user where id = \? order by email asc limit 1000$`).WithArgs(1).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("*").Where("id = ?", 1).OrderBy("email asc").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyOffset(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"id", "role", "email", "active"})
	for i := range control {
		rows.AddRow(control[i].ID, control[i].Role, control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select \* from user limit 1000 offset 5$`).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("*").Offset(5).To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectManyJoinSomeFields(t *testing.T) {
	control := makeUserSlice(3)
	rows := sqlmock.NewRows([]string{"email", "active"})
	for i := range control {
		rows.AddRow(control[i].Email, control[i].Active)
	}
	mock.ExpectQuery(`^select user\.email, user\.active from user join post on post.user_id = user.id limit 1000$`).WillReturnRows(rows)

	var test []*User
	if err := DBConn.Select("user.email", "user.active").Join("post on post.user_id = user.id").To(&test); err != nil {
		t.Errorf("error was not expected while selecting: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	for i := range control {
		test[i].ID = control[i].ID
		test[i].Role = control[i].Role
	}
	if err := assertSameSlice(control, test); err != nil {
		t.Error(err)
	}
}

func TestBuilderSelectOneNoResult(t *testing.T) {
	mock.ExpectQuery(`^select \* from user where id = \? limit 1$`).WithArgs(1).WillReturnRows(mock.NewRows([]string{}))

	if err := DBConn.Select("*").Where("id = ?", 1).To(&User{}); err != ErrNotFound {
		t.Errorf("expected ErrNotFound: got %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestBuilderSelectOneErrors(t *testing.T) {
	assertErr(
		t,
		"should return error if non pointer passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(User{}),
	)

	type Post struct {
		ID    int64
		Title string
	}
	assertErr(
		t,
		"should return error if unregistered struct passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(&Post{}),
	)
	assertErr(
		t,
		"should return error if slice of non pointers passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(&[]Post{}),
	)

	type testMap map[string]interface{}
	assertErr(
		t,
		"should return error if slice of pointers to non structs passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(&[]*testMap{}),
	)

	assertErr(
		t,
		"should return error if slice of pointers to non registered items passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(&[]*Post{}),
	)

	testStr := "asdf"
	assertErr(
		t,
		"should return error if non struct and non slice passed to To",
		DBConn.Select("*").Where("id = ?", 1).To(&testStr),
	)
}
