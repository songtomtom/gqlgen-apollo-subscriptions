import React, { ChangeEvent, useEffect, useState } from 'react';
import { gql, useMutation, useQuery } from '@apollo/client';

import './App.css';

const CREATE_POST = gql`
  mutation CreatePost($input: CreatePostInput!) {
    createPost(input: $input) {
      id
    }
  }
`;

const CREATE_COMMENT = gql`
  mutation CreateComment($input: CreateCommentInput!) {
    createComment(input: $input) {
      id
      postId
      content
    }
  }
`;

const LIST_COMMENTS = gql`
  query ListComments($where: CommentsWhere!) {
    comments(where: $where) {
      id
      postId
      content
    }
  }
`;

const COMMENTS_SUBSCRIPTION = gql`
  subscription OnCommentAdded($input: AddedCommentInput!) {
    commentAdded(input: $input) {
      id
      postId
      content
    }
  }
`;

function App() {
  const [id, setId] = useState<string>('songtomtom');
  const [content, setContent] = useState<string>('Hello~!');

  const [createPost] = useMutation(CREATE_POST, {
    variables: { input: { id } },
  });
  const [createComment] = useMutation(CREATE_COMMENT, {
    variables: { input: { postId: id, content } },
  });
  // subscribeToMore 추가
  const { data, loading, subscribeToMore } = useQuery(LIST_COMMENTS, {
    variables: { where: { postId: id } },
  });

  const onClickCreatePost = async () => {
    const { data } = await createPost();
    if (data) {
      // Do something
    }
  };

  const onClickCreateComment = async () => {
    const { data } = await createComment();
    if (data) {
      // Do something
    }
  };

  const onChangeId = (e: ChangeEvent<HTMLInputElement>) => {
    setId(e.target.value);
  };

  const onChangeContent = (e: ChangeEvent<HTMLInputElement>) => {
    setContent(e.target.value);
  };

  const subscribeToNewComment = () => {
    return subscribeToMore({
      document: COMMENTS_SUBSCRIPTION,
      variables: {
        input: {
          postId: id,
        },
      },
      updateQuery: (prev, { subscriptionData }) => {
        if (!subscriptionData.data) {
          return prev;
        }
        const {
          data: { commentAdded: newComment },
        } = subscriptionData;
        return Object.assign({}, prev, {
          comments: [newComment, ...prev.comments],
        });
      },
    });
  };

  useEffect(() => subscribeToNewComment(), []);

  return (
    <div>
      <div>
        <input type="text" onChange={onChangeId} value={id} />
        id
        <br />
        <input type="text" onChange={onChangeContent} value={content} />
        content
      </div>

      <button onClick={onClickCreatePost}>create post</button>
      <button onClick={onClickCreateComment}>create comment</button>
      <ul>
        {!loading &&
          data.comments.map((item: any, index: number) => {
            return (
              <li key={index}>
                <small>{item.postId}</small>: <strong>{item.content}</strong>
              </li>
            );
          })}
      </ul>
    </div>
  );
}

export default App;
