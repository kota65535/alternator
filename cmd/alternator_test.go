package cmd

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestFetchGlobalConfig(t *testing.T) {
	uri := fmt.Sprintf("mysql://root@localhost:13306/?multiStatements=true")
	dbUri, err := NewDatabaseUri(uri)
	assert.NoError(t, err)

	db, err := sql.Open(dbUri.Dialect, dbUri.Dsn())
	assert.NoError(t, err)

	config, err := fetchGlobalConfig(db)
	assert.NoError(t, err)

	assert.NotEmpty(t, config.CharacterSetDatabase)
	assert.NotEmpty(t, config.CharacterSetServer)
	assert.NotEmpty(t, config.CollationServer)
	assert.NotEmpty(t, config.CharsetToCollation)
}
