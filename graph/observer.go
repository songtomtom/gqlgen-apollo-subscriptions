package graph

import (
	"sync"

	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
)

// Observer 는 postId 별 구독자 채널을 관리한다.
//
// 구독(Subscription) 리졸버는 Subscribe 로 채널을 받아 gqlgen 에 돌려주고,
// 뮤테이션 리졸버는 Publish 로 같은 postId 를 구독 중인 모든 채널에 이벤트를 보낸다.
//
// 설계 메모
//   - postId 하나에 여러 구독자가 붙을 수 있어야 하므로 map[postId]map[chan]struct{} 구조를 쓴다.
//   - 리졸버는 요청마다 다른 고루틴에서 실행되므로 map 접근은 RWMutex 로 보호한다.
//   - Publish 는 절대 블로킹되면 안 된다. 느린 구독자 하나 때문에 뮤테이션 전체가 멈추면 안 되기 때문에
//     채널에 버퍼를 두고, 그래도 가득 차 있으면 그 구독자에 대한 이벤트는 버린다.
type Observer struct {
	mu   sync.RWMutex
	subs map[string]map[chan *model.Comment]struct{}
}

// subscriberBuffer 는 구독자 한 명이 소비하지 못한 채 쌓아 둘 수 있는 이벤트 수다.
const subscriberBuffer = 16

func NewObserver() *Observer {
	return &Observer{subs: map[string]map[chan *model.Comment]struct{}{}}
}

// Subscribe 는 postID 에 대한 수신 채널과 구독 해제 함수를 돌려준다.
// 구독 해제 함수는 채널을 등록에서 제거하고 닫는다. 두 번 호출해도 안전하다.
func (o *Observer) Subscribe(postID string) (<-chan *model.Comment, func()) {
	ch := make(chan *model.Comment, subscriberBuffer)

	o.mu.Lock()
	if o.subs[postID] == nil {
		o.subs[postID] = map[chan *model.Comment]struct{}{}
	}
	o.subs[postID][ch] = struct{}{}
	o.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			o.mu.Lock()
			defer o.mu.Unlock()
			if set, ok := o.subs[postID]; ok {
				delete(set, ch)
				if len(set) == 0 {
					delete(o.subs, postID)
				}
			}
			close(ch)
		})
	}
	return ch, unsubscribe
}

// Publish 는 postID 를 구독 중인 모든 채널에 comment 를 보낸다.
// 버퍼가 가득 찬 구독자는 건너뛰며, 실제로 전달된 구독자 수를 돌려준다.
func (o *Observer) Publish(postID string, comment *model.Comment) int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	delivered := 0
	for ch := range o.subs[postID] {
		select {
		case ch <- comment:
			delivered++
		default:
			// 구독자가 이벤트를 소비하지 못하고 있다. 뮤테이션을 막지 않기 위해 버린다.
		}
	}
	return delivered
}

// Count 는 postID 의 현재 구독자 수를 돌려준다. 테스트와 디버깅용.
func (o *Observer) Count(postID string) int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.subs[postID])
}
