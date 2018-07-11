package object

import (
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud-message/slot"
	"testing"
)

type DefinitionSuite struct {
	suite.Suite
	helmet *Definition
}

func TestDefinitionSuite(t *testing.T) {
	suite.Run(t, new(DefinitionSuite))
}

func (suite *DefinitionSuite) SetupTest() {
	suite.helmet = NewDefinition("definitionId", "helmet", "zoneId",
		Armor, []string{"iron", "helm"}, "desc", "desc on ground", slot.Head)
}

func (suite *DefinitionSuite) TestHasAlias() {
	suite.Assert().True(suite.helmet.HasAlias("helm"), "should have alias")
	suite.Assert().False(suite.helmet.HasAlias("bronze"), "should not have alias")
}
