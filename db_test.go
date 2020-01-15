package gosql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/twharmon/gosql"
)

func TestDelete(t *testing.T) {
	type DeleteModel struct {
		ID int `gosql:"primary"`
	}
	check(t, gosql.Register(DeleteModel{}))
	deleteModel := DeleteModel{5}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^delete from delete_model where id = \?$`).WithArgs(deleteModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Delete(&deleteModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	type UpdateModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(UpdateModel{}))
	updateModel := UpdateModel{5, "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^update update_model set name = \? where id = \?$`).WithArgs(updateModel.Name, updateModel.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Update(&updateModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsert(t *testing.T) {
	type InsertModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(InsertModel{}))
	insertModel := InsertModel{Name: "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^insert into insert_model \(name\) values \(\?\)$`).WithArgs(insertModel.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModel)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}

func TestInsertWithPrimary(t *testing.T) {
	type InsertWithPrimaryModel struct {
		ID   int `gosql:"primary"`
		Name string
	}
	check(t, gosql.Register(InsertWithPrimaryModel{}))
	insertModelWithPrimary := InsertWithPrimaryModel{5, "foo"}
	db, mock, err := getMockDB()
	check(t, err)
	mock.ExpectExec(`^insert into insert_with_primary_model \(id, name\) values \(\?, \?\)$`).WithArgs(insertModelWithPrimary.ID, insertModelWithPrimary.Name).WillReturnResult(sqlmock.NewResult(0, 1))
	_, err = db.Insert(&insertModelWithPrimary)
	check(t, err)
	check(t, mock.ExpectationsWereMet())
}
