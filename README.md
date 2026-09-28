# gqlgen + Apollo Client 구독(Subscription) 예제

Go [gqlgen](https://gqlgen.com) 으로 GraphQL Subscription 서버를 만들고, React + [Apollo Client](https://www.apollographql.com/docs/react/) 에서 WebSocket 으로 구독해 댓글이 달릴 때마다 목록이 실시간으로 갱신되는 예제입니다.

블로그 글 5편과 함께 읽도록 만들었습니다.

1. [gqlgen 으로 구독(Subscriptions) 서버 만들기](https://songtomtom.github.io/blog/gqlgen-subscriptions-server)
2. [Apollo Client 클라이언트 WebSocket Link 연결](https://songtomtom.github.io/blog/apollo-client-websocket-link)
3. [실시간 구독 서비스 리졸버 만들기](https://songtomtom.github.io/blog/realtime-subscription-resolver)
4. [Apollo Client로 구독 서비스 GUI 만들기](https://songtomtom.github.io/blog/apollo-client-subscription-gui)
5. [Channel Observer를 사용하여 구독 서비스 관리](https://songtomtom.github.io/blog/channel-observer-subscription)

## 구조

```
클라이언트 (React + Apollo Client)
  ├─ HttpLink  ──POST──▶  /query          Query, Mutation
  └─ WsLink    ──WS────▶  /subscriptions  Subscription (graphql-transport-ws)

서버 (Go + gqlgen)
  createComment 뮤테이션 ──▶ MySQL 저장 ──▶ Observer.Publish(postId)
                                                  │
  commentAdded 구독 ◀── Observer.Subscribe(postId) ┘  (postId 별 구독자 채널 집합)
```

- `server.go`: HTTP 핸들러와 WebSocket 전송 설정
- `graph/schema.graphqls`: 스키마
- `graph/schema.resolvers.go`: 리졸버
- `graph/observer.go`: postId 별 구독자 채널을 관리하는 Observer (mutex, 논블로킹 발행, 구독 해제 정리)
- `client/`: Create React App + Apollo Client

## 실행

```bash
make up        # MySQL 컨테이너 시작 (docker compose)
make start     # 서버 시작, http://localhost:8080 에 playground
make client    # 클라이언트 시작, http://localhost:3000
```

환경 변수 `DSN`, `PORT` 로 MySQL 접속 정보와 포트를 바꿀 수 있습니다.

스키마를 바꾼 뒤에는 `make gen` 으로 리졸버 코드를 다시 생성합니다.

## 테스트

```bash
go test ./...
```

Observer 의 동시성, 구독 해제, 논블로킹 발행을 검증하는 단위 테스트가 들어 있습니다.
