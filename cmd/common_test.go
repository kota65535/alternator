package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestFetchGlobalConfig(t *testing.T) {
	dsn := fmt.Sprintf("mysql://root@localhost:13306/?multiStatements=true")
	Db = ConnectToDb(ParseDatabaseUrl(dsn))
	defer Db.Close()
	config := FetchGlobalConfig()

	assert.NotEmpty(t, config.CharacterSetDatabase)
	assert.NotEmpty(t, config.CharacterSetServer)
	assert.NotEmpty(t, config.CollationServer)
	assert.NotEmpty(t, config.CharsetToCollation)
}
