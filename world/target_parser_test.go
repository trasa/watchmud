package world

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TargetParserSuite struct {
	suite.Suite
}

func TestTargetParserSuite(t *testing.T) {
	suite.Run(t, new(TargetParserSuite))
}

func (suite *TargetParserSuite) TestParse() {
	target, err := parseTarget("foo")
	suite.Assert().NoError(err)
	suite.Assert().Equal("foo", target.Name)
	suite.Assert().Equal(0, target.Identifier)
	suite.Assert().Equal(0, target.Quantity)
}

func (suite *TargetParserSuite) TestParseAll() {
	target, err := parseTarget("all")
	suite.Assert().NoError(err)
	suite.Assert().True(target.All)
	suite.Assert().Equal("", target.Name)
	suite.Assert().Equal(0, target.Identifier)
	suite.Assert().Equal(0, target.Quantity)
}

func (suite *TargetParserSuite) TestParseWithIdentifier() {
	target, err := parseTarget("2.foo")
	suite.Assert().NoError(err)
	suite.Assert().Equal("foo", target.Name)
	suite.Assert().Equal(2, target.Identifier)
	suite.Assert().Equal(0, target.Quantity)
}
func (suite *TargetParserSuite) TestParseWithAllIdentifier() {
	target, err := parseTarget("all.foo")
	suite.Assert().NoError(err)
	suite.Assert().Equal("foo", target.Name)
	suite.Assert().True(target.All)
	suite.Assert().Equal(0, target.Identifier)
	suite.Assert().Equal(0, target.Quantity)
}
func (suite *TargetParserSuite) TestParseQuantity() {
	target, err := parseTarget("50 foo")
	suite.Assert().NoError(err)
	suite.Assert().Equal("foo", target.Name)
	suite.Assert().Equal(0, target.Identifier)
	suite.Assert().Equal(50, target.Quantity)
}

func (suite *TargetParserSuite) TestParseQuantityNotNumber() {
	_, err := parseTarget("x foo")
	suite.Assert().Error(err)
}

func (suite *TargetParserSuite) TestParseQuantityWithIdentifier() {
	target, err := parseTarget("50 2.foo")
	suite.Assert().NoError(err)
	suite.Assert().Equal("foo", target.Name)
	suite.Assert().Equal(2, target.Identifier)
	suite.Assert().Equal(50, target.Quantity)
}

func (suite *TargetParserSuite) TestParseEmpty() {
	target, err := parseTarget("")
	suite.Assert().NoError(err)
	suite.Assert().Equal("", target.Name)
	suite.Assert().Equal(0, target.Identifier)
	suite.Assert().Equal(0, target.Quantity)
}

func (suite *TargetParserSuite) TestParseTooManyParts() {
	_, err := parseTarget("foo bar baz")
	suite.Assert().Error(err)
}

func (suite *TargetParserSuite) TestParseTooManyDots() {
	_, err := parseTarget("50 5.bar.baz")
	suite.Assert().Error(err)
}
