package telnet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/rules"
)

// a two-species catalog with deliberately colliding first letters: Hill Dwarf
// and High Elf both start with "h", which is the case the prefix rule exists
// for.
func testLineageCatalog(t *testing.T) *rules.Catalog {
	t.Helper()
	dwarf := &rules.Species{Id: "dwarf", Name: "Dwarf", Lineages: []*rules.Lineage{
		{Id: "hill_dwarf", Name: "Hill Dwarf", Description: "from the terraces"},
		{Id: "mountain_dwarf", Name: "Mountain Dwarf"},
	}}
	elf := &rules.Species{Id: "elf", Name: "Elf", Lineages: []*rules.Lineage{
		{Id: "high_elf", Name: "High Elf"},
	}}
	cat, err := rules.NewCatalog([]*rules.Species{dwarf, elf}, rules.NewTestRoles())
	require.NoError(t, err)
	return cat
}

func TestMatchLineage(t *testing.T) {
	choices := lineageChoices(testLineageCatalog(t))
	require.Len(t, choices, 3)

	cases := []struct {
		name   string
		answer string
		want   string // empty means "no match, ask again"
	}{
		{"by menu number", "2", "mountain_dwarf"},
		{"first number", "1", "hill_dwarf"},
		{"last number", "3", "high_elf"},
		{"number below the menu", "0", ""},
		{"number past the menu", "4", ""},
		{"a number that isn't one", "-1", ""},
		{"by full name", "Mountain Dwarf", "mountain_dwarf"},
		{"by name, any case", "mountain dwarf", "mountain_dwarf"},
		{"by id", "high_elf", "high_elf"},
		{"by unambiguous prefix", "moun", "mountain_dwarf"},
		// "h" fits both Hill Dwarf and High Elf: guessing one would quietly
		// hand the player a character they didn't ask for
		{"ambiguous prefix", "h", ""},
		// "hi" still fits both Hill Dwarf and High Elf; "hil" doesn't
		{"ambiguous two letters", "hi", ""},
		{"prefix long enough to disambiguate", "hil", "hill_dwarf"},
		{"surrounding whitespace is trimmed", "  high elf  ", "high_elf"},
		{"nothing like it", "gnome", ""},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, found := matchLineage(choices, tc.answer)
			if tc.want == "" {
				assert.False(t, found, "wanted no match, got %q", got)
				return
			}
			require.True(t, found)
			assert.Equal(t, tc.want, got)
		})
	}
}

// The menu is what the player picks numbers off, so the numbering has to run
// straight through the species groups rather than restarting in each.
func TestLineageMenuNumbersRunThroughTheGroups(t *testing.T) {
	menu := lineageMenu(testLineageCatalog(t))
	for _, want := range []string{"Dwarf", " 1) Hill Dwarf", " 2) Mountain Dwarf", "Elf", " 3) High Elf"} {
		assert.Contains(t, menu, want)
	}
	// descriptions ride along where content gives one
	assert.Contains(t, menu, "from the terraces")
	// and the menu says what the choice does and doesn't do
	assert.Contains(t, menu, "It decides nothing but how you look")
	assert.True(t, strings.HasPrefix(menu, "\n"))
}
