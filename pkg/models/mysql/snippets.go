package mysql

import (
	"alex/pkg/models"
	"database/sql"
)

// Define a Snippet type which wraps a sql.Db connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	/*Write the SQL statement we want to execute. I've split it over two lines
	for readability (which is why it's surrounded with backquotes instead
	of normal double quotes)*/
	stmt := `INSERT INTO snippets(title,content,created,expires)
	VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	/*Use the Exec()method on the embeded connection pool to execute the statement
	The first parametr is the Sql statement, followed by the title,content,expiry values for the placeholder parametrs
	This method returns a sql.Result object, which containssome basic
	information about what happened when the statement was executed.
	*/
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	/*Use the LastInsertId() method on the result object to get the ID of our
	newly inserted record in the snippets table*/
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The Id returned has the type int64, so  we convert it to an int type
	//before returnning
	return int(id), nil
}

// This will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	/*Write the SQL statement we want to execute. Again, I've split it two
	lines for readability
	*/
	stmt := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	//Use thr QueryRow() method on the connection pool to execute our
	//Sql statement, passing in the untrusted id variable as the value for the
	//placholder parametr. This returns a pointer to a sql.Row object which
	//holds the result from the database.
	row := m.DB.QueryRow(stmt, id) // Здесь храниться результат из БД
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	//Coresponding field in the Snippet struct.Notice that the arguments
	//to row.Scan are *pointers* to the place you want to copy the data info,
	//and the number of arguments must be exactly the same as the number of columns
	// returned by your statement.
	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows { // If the Query returns no rows, then row scan() will return
		// our own models.ErrNoRecord error instead of a Snippet object
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// Write a Sql statement we want to execute.
	smt := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	//Use the Query() method on the connection pool to execute our Sql statement
	//This returns a sql.Rows resultset containing the result our query
	rows, err := m.DB.Query(smt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next() { // Если итерация по всем строкам завершена, то набор результатов
		// автоматически закрывается и освобождает базовое соединения с базой данных
		s := &models.Snippet{}
		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// When the rows.Next()loop has finished we call rows.Err() to retrieve any
	//error that was encountered during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went ok then return the Snippets slice
	return snippets, nil
}
