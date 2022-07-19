package parser

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestFetchConfig(t *testing.T) {
	dsn := fmt.Sprintf("root@(localhost:3306)/?multiStatements=true")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	config := FetchGlobalConfig(db)

	assert.NotEmpty(t, config.CharacterSetDatabase)
	assert.NotEmpty(t, config.CharacterSetServer)
	assert.NotEmpty(t, config.CollationServer)
	assert.NotEmpty(t, config.CharsetToCollation)
}