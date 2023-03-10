package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"time"

	"github.com/dgryski/trifles/uuid"
	"github.com/songtomtom/gqlgen-apollo-subscriptions/graph/model"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	post := model.Post{
		ID: input.ID,
	}
	r.DB.Create(post)

	return &post, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	comment := model.Comment{
		ID:      uuid.UUIDv4(),
		PostID:  input.PostID,
		Content: input.Content,
	}
	r.DB.Create(comment)

	for _, o := range r.Observer {
		o <- &comment
	}

	return &comment, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, where model.CommentsWhere) ([]*model.Comment, error) {
	_, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var comments []*model.Comment
	r.DB.Where("post_id = ?", where.PostID).Find(&comments)

	return comments, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, input model.AddedCommentInput) (<-chan *model.Comment, error) {

	ch := make(chan *model.Comment)

	go func() {
		<-ctx.Done()
		delete(r.Observer, input.PostID)
	}()

	r.Observer[input.PostID] = ch

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
