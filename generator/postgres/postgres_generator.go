package postgres

import (
	"database/sql"
	"fmt"
	"path"
	"strings"

	"github.com/go-jet/jet/generator/internal/metadata"
	"github.com/go-jet/jet/generator/internal/template"
	"github.com/go-jet/jet/internal/utils"
	"github.com/go-jet/jet/postgres"
)

// DBConnection contains postgres connection details
type DBConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	SslMode  string
	Params   string

	DBName     string
	SchemaName string
}

// Generate generates jet files at destination dir from database connection details
func Generate(destDir string, dbConn DBConnection) (err error) {
	defer utils.ErrorCatch(&err)

	db, err := openConnection(dbConn)
	utils.PanicOnError(err)
	defer utils.DBClose(db)

	fmt.Println("Retrieving schema information...")
	schemaInfo := metadata.GetSchemaMetaData(db, dbConn.SchemaName, &postgresQuerySet{})

	genPath := path.Join(destDir, dbConn.DBName, dbConn.SchemaName)
	template.GenerateFiles(genPath, schemaInfo, postgres.Dialect)

	return
}

// escapePostgresValue escapes a connection string parameter's value.
//
// For details, see section 33.1.1.1 of
// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func escapePostgresValue(in string) string {
	in = strings.ReplaceAll(in, `\`, `\\`)
	in = strings.ReplaceAll(in, `'`, `\'`)
	return in
}

func openConnection(dbConn DBConnection) (*sql.DB, error) {
	// By only adding the clauses specified in dbConn (as indicated by non-zero
	// values), lower level postgresql libraries retain the ability to discover
	// values set via environment variables, as per the postgresql docs.  For
	// more information, see session 33.14 of
	// https://www.postgresql.org/docs/current/libpq-envars.html
	clauses := []string{}
	if dbConn.Host != "" {
		clauses = append(clauses, fmt.Sprintf("host='%s'", dbConn.Host))
	}
	if dbConn.Port > 0 {
		clauses = append(clauses, fmt.Sprintf("port=%d", dbConn.Port))
	}
	if dbConn.User != "" {
		clauses = append(clauses, fmt.Sprintf("user='%s'", dbConn.User))
	}
	if dbConn.Password != "" {
		clauses = append(clauses, fmt.Sprintf("password='%s'", dbConn.Password))
	}
	if dbConn.DBName != "" {
		clauses = append(clauses, fmt.Sprintf("dbname='%s'", dbConn.DBName))
	}
	if dbConn.SslMode != "" {
		clauses = append(clauses, fmt.Sprintf("sslmode='%s'", dbConn.SslMode))
	}
	connectionString := fmt.Sprintf("%s %s",
		strings.Join(clauses, " "), dbConn.Params)

	fmt.Println("Connecting to postgres database: " + connectionString)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
