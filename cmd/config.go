package cmd

import (
	"database/sql"
	"github.com/kota65535/alternator/parser"
	"github.com/spf13/cobra"
)

func fetchGlobalConfig(db *sql.DB) *parser.GlobalConfig {
	rows1, err := db.Query("SHOW GLOBAL VARIABLES")
	defer rows1.Close()
	cobra.CheckErr(err)
	var name string
	var value string
	variables := map[string]string{}
	for rows1.Next() {
		err = rows1.Scan(&name, &value)
		variables[name] = value
	}

	rows2, err := db.Query("SHOW CHARACTER SET")
	defer rows2.Close()
	cobra.CheckErr(err)
	var charset string
	var description string
	var collation string
	var maxLen string
	charsetToCollation := map[string]string{}
	for rows2.Next() {
		err = rows2.Scan(&charset, &description, &collation, &maxLen)
		charsetToCollation[charset] = collation
	}

	return &parser.GlobalConfig{
		CharacterSetServer:   variables["character_set_server"],
		CharacterSetDatabase: variables["character_set_database"],
		CollationServer:      variables["collation_server"],
		CharsetToCollation:   charsetToCollation,
	}
}
