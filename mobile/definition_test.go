package mobile

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/watchmud/watchmud/rules"
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
		rules.WanderDefinition{CanWander: false},
		10,
		false)
}

func (suite *DefinitionSuite) TestFlags() {
	suite.Assert().False(suite.definition.HasFlag("blah"))
	suite.definition.flags["blah"] = true
	suite.Assert().True(suite.definition.HasFlag("blah"))
}

func (suite *DefinitionSuite) TestSetFlags() {
	suite.definition.SetFlags([]Flag{"Aggressive", "PlayerCantFight"})
	suite.definition.SetFlags(nil)
	suite.definition.SetFlags([]Flag{})

	suite.Assert().True(suite.definition.HasFlag("Aggressive"))

}
