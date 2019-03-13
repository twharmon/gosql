package gosql

// Exec .
func (db *DB) Exec(query string, args ...interface{}) error {
	_, err := db.db.Exec(query, args...)
	return err
}
