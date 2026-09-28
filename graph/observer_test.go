package graph

import (
	"sync"
	"testing"

	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
)

func TestObserver_PublishReachesAllSubscribersOfSamePost(t *testing.T) {
	o := NewObserver()
	a, unsubA := o.Subscribe("post-1")
	b, unsubB := o.Subscribe("post-1")
	other, unsubOther := o.Subscribe("post-2")
	defer unsubA()
	defer unsubB()
	defer unsubOther()

	c := &model.Comment{ID: "c1", PostID: "post-1", Content: "hi"}
	if got := o.Publish("post-1", c); got != 2 {
		t.Fatalf("delivered = %d, want 2", got)
	}
	if got := <-a; got != c {
		t.Fatalf("subscriber a got %v", got)
	}
	if got := <-b; got != c {
		t.Fatalf("subscriber b got %v", got)
	}
	select {
	case got := <-other:
		t.Fatalf("post-2 subscriber should not receive post-1 comment, got %v", got)
	default:
	}
}

func TestObserver_UnsubscribeRemovesOnlyItself(t *testing.T) {
	o := NewObserver()
	_, unsubA := o.Subscribe("post-1")
	b, unsubB := o.Subscribe("post-1")
	defer unsubB()

	unsubA()
	unsubA() // 두 번 호출해도 panic 이 없어야 한다

	if got := o.Count("post-1"); got != 1 {
		t.Fatalf("count = %d, want 1", got)
	}
	if got := o.Publish("post-1", &model.Comment{ID: "c"}); got != 1 {
		t.Fatalf("delivered = %d, want 1", got)
	}
	<-b
}

func TestObserver_PublishDoesNotBlockOnSlowSubscriber(t *testing.T) {
	o := NewObserver()
	_, unsub := o.Subscribe("post-1")
	defer unsub()

	// 아무도 읽지 않는 채널에 버퍼보다 많이 보내도 Publish 는 돌아와야 한다.
	for i := 0; i < subscriberBuffer*2; i++ {
		o.Publish("post-1", &model.Comment{ID: "c"})
	}
}

func TestObserver_ConcurrentAccess(t *testing.T) {
	o := NewObserver()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, unsub := o.Subscribe("post-1")
			unsub()
		}()
		go func() {
			defer wg.Done()
			o.Publish("post-1", &model.Comment{ID: "c"})
		}()
	}
	wg.Wait()
	if got := o.Count("post-1"); got != 0 {
		t.Fatalf("count = %d, want 0", got)
	}
}
