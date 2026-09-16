package object

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/rules"
)

type DefinitionSuite struct {
	suite.Suite
	helmet *Definition
}

func TestDefinitionSuite(t *testing.T) {
	suite.Run(t, new(DefinitionSuite))
}

func (suite *DefinitionSuite) SetupTest() {
	suite.helmet = NewDefinition(
		"definitionId",
		"helmet",
		"zoneId",
		Armor,
		[]string{"iron", "helm"},
		"desc",
		"desc on ground",
		rules.SlotHead,
		"plate", // TODO
	)
}

func (suite *DefinitionSuite) TestHasAlias() {
	suite.Assert().True(suite.helmet.HasAlias("helm"), "should have alias")
	suite.Assert().False(suite.helmet.HasAlias("bronze"), "should not have alias")
}
