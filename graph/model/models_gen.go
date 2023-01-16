// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

type CommentsWhere struct {
	// 유저 프로필 아이디
	PostID string `json:"postId"`
}

type CreateCommentInput struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

type Post struct {
	ID string `json:"id"`
}
