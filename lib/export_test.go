package lib

import (
	"fmt"
	"github.com/kota65535/alternator/parser"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// Define common functions only used by test files here

var TestDefaultGlobalConfig = &parser.GlobalConfig{
	CharacterSetServer: "utf8mb4",
	CollationServer:    "utf8mb4_0900_ai_ci",
	CharsetToCollation: map[string]string{
		"utf8mb4": "utf8mb4_0900_ai_ci",
		"utf16":   "utf16_general_ci",
	},
}

func getAlteredDatabases(t *testing.T, q1 string, q2 string) DatabaseAlterations {
	f1, err := os.Open(q1)
	require.NoError(t, err)

	p1 := parser.NewParser(f1)
	r1, err := p1.Parse()
	require.NoError(t, err)

	f2, err := os.Open(q2)
	require.NoError(t, err)

	p2 := parser.NewParser(f2)
	r2, err := p2.Parse()
	require.NoError(t, err)

	s1 := normalizeStatements(r1, TestDefaultGlobalConfig)
	s2 := normalizeStatements(r2, TestDefaultGlobalConfig)

	fmt.Println("========== Tables ==========")
	for _, s := range s1[0].Tables {
		fmt.Println(s.String())
	}
	fmt.Println("====================")
	for _, s := range s2[0].Tables {
		fmt.Println(s.String())
	}

	alt := NewDatabaseAlterations(s1, s2)
	fmt.Println("========== Statements ==========")
	for _, s := range alt.Statements() {
		fmt.Println(s)
	}
	fmt.Println("=====================")
	return alt
}
