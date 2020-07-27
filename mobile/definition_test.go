package mobile

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DefinitionSuite struct {
	suite.Suite
	definition *Definition
}

func TestDefinitionSuite(t *testing.T) {
	suite.Run(t, new(DefinitionSuite))
}

func (suite *DefinitionSuite) SetupTest() {
	suite.definition = NewDefinition("definitionId",
		"name",
		"zone",
		[]string{"alias"},
		"short desc",
		"descr",
		25,
		WanderingDefinition{CanWander: false},
		10)
}

func (suite *DefinitionSuite) TestFlags() {
	suite.Assert().False(suite.definition.HasFlag("blah"))
	suite.definition.flags["blah"] = true
	suite.Assert().True(suite.definition.HasFlag("blah"))
}

func (suite *DefinitionSuite) TestSetFlags() {
	suite.definition.SetFlags([]string{"a", "b"})
	suite.definition.SetFlags(nil)
	suite.definition.SetFlags([]string{})

	suite.Assert().True(suite.definition.HasFlag("a"))

}
