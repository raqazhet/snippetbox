package mysql

import "database/sql"

type ExampleModel struct {
	DB *sql.DB
}

func (m *ExampleModel) ExampleTransaction() error {
	// Calling the Begin()method on the connection a new sql.Tx
	tx, err := m.DB.Begin() // this object represents the in-progress database transaction
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		// if there is any error, we call the tx.Rollback method
		// This will abort the transaction and no changes
		tx.Rollback() //
		return err
	}
	_, err = tx.Exec("Update ...")
	if err != nil {
		tx.Rollback()
		return err
	}
	// if there no errors, the statements in the transaction can be commited
	// to the database with the tx.Commit() method
	err = tx.Commit()
	return err

}
