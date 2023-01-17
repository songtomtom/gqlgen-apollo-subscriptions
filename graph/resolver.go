package graph

import (
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB       *gorm.DB
	Observer map[string]chan *model.Comment // 추가
}
