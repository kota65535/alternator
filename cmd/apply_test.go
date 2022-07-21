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

var ApplyParam = ApplyParams{
	AutoApprove: true,
}

var ApplyTestFixtures = []Fixture{}

func init() {
	dirs := getDirs(RootPath)
	sort.Strings(dirs)
	for _, dir := range dirs {
		for _, db := range Databases {
			ApplyTestFixtures = append(ApplyTestFixtures, Fixture{db, dir})
		}
	}
	fmt.Println(ApplyTestFixtures)
}

func TestApply(t *testing.T) {
	for _, fixture := range ApplyTestFixtures {
		t.Run(fixture.Name(), func(t *testing.T) {
			if lib.Contains(Skipped, fixture.Name()) {
				t.Skip()
			}
			url := fmt.Sprintf("%s://root@localhost:%d/", fixture.Dialect, fixture.Port)
			dir := fixture.Dir

			// when
			err := prepareDb(dir, fixture.Database)
			require.NoError(t, err)
			alt := ApplyCmd(testFile(dir, "to.sql", fixture.Database), url, ApplyParam)

			// assert ALTER statements
			s, err := getAlter(dir, fixture.Database)
			require.NoError(t, err)
			assert.Equal(t, s, strings.Join(alt.Statements(), "\n"))

			// 2nd apply should return empty statements
			alt2 := ApplyCmd(testFile(dir, "to.sql", fixture.Database), url, ApplyParam)
			assert.Nil(t, alt2)
		})
	}
}
