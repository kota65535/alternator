package cmd

import (
	"fmt"
	"github.com/kota65535/alternator/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

var DiffTestFixtures = []Fixture{}

func init() {
	dirs := getDirs(RootPath)
	sort.Strings(dirs)
	for _, dir := range dirs {
		for _, db := range Databases {
			DiffTestFixtures = append(DiffTestFixtures, Fixture{db, dir})
		}
	}
	fmt.Println(DiffTestFixtures)
}

func TestDiff(t *testing.T) {
	for _, fixture := range DiffTestFixtures {
		t.Run(fixture.Name(), func(t *testing.T) {
			if lib.Contains(Skipped, fixture.Name()) {
				t.Skip()
			}
			url := fmt.Sprintf("%s://root@localhost:%d/", fixture.Dialect, fixture.Port)
			dir := fixture.Dir

			// when
			err := prepareDb(dir, fixture.Database)
			require.NoError(t, err)
			alt := diffCmd(testFile(dir, "to.sql", fixture.Database), url)

			// assert diff
			s, err := getDiff(dir, fixture.Database)
			require.NoError(t, err)
			assert.Equal(t, s, strings.Join(alt.Diff(), "\n"))
		})
	}
}
