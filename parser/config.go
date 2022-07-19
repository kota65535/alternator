package parser

import (
	"database/sql"
	"github.com/spf13/cobra"
)

type GlobalConfig struct {
	CharacterSetServer   string
	CharacterSetDatabase string
	CollationServer      string
	CharsetToCollation   map[string]string
}

func FetchGlobalConfig(db *sql.DB) *GlobalConfig {
	rows, err := db.Query("SHOW GLOBAL VARIABLES")
	defer rows.Close()
	cobra.CheckErr(err)
	var name string
	var value string
	variables := map[string]string{}
	for rows.Next() {
		err = rows.Scan(&name, &value)
		variables[name] = value
	}

	rows, err = db.Query("SHOW CHARACTER SET")
	defer rows.Close()
	cobra.CheckErr(err)
	var charset string
	var description string
	var collation string
	var maxLen string
	charsetToCollation := map[string]string{}
	for rows.Next() {
		err = rows.Scan(&charset, &description, &collation, &maxLen)
		charsetToCollation[charset] = collation
	}

	return &GlobalConfig{
		CharacterSetServer:   variables["character_set_server"],
		CharacterSetDatabase: variables["character_set_database"],
		CollationServer:      variables["collation_server"],
		CharsetToCollation:   charsetToCollation,
	}
}
