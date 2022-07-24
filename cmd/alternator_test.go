package cmd

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestFetchGlobalConfig(t *testing.T) {
	uri := fmt.Sprintf("mysql://root@localhost:13306/mydb")
	dbUri, err := NewDatabaseUri(uri)
	require.NoError(t, err)

	db, err := sql.Open(dbUri.Dialect, dbUri.DsnWoDbName())
	require.NoError(t, err)

	config, err := fetchGlobalConfig(db)
	require.NoError(t, err)

	assert.NotEmpty(t, config.CharacterSetDatabase)
	assert.NotEmpty(t, config.CharacterSetServer)
	assert.NotEmpty(t, config.CollationServer)
	assert.NotEmpty(t, config.CharsetToCollation)
}
