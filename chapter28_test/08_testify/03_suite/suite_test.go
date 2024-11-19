package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestExample(t *testing.T) {
	suite.Run(t, new(MyTestSuit))
}
