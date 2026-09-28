package graph

import (
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver 는 모든 리졸버가 공유하는 의존성이다.
// 서버 시작 시 한 번 만들어지고, 요청마다 새로 만들지 않는다.
// Observer 를 요청마다 새로 만들면 구독자와 발행자가 서로 다른 Observer 를 보게 되어 이벤트가 전달되지 않는다.
type Resolver struct {
	DB       *gorm.DB
	Observer *Observer
}
