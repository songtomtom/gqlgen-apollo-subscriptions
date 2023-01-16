import React, { ChangeEvent, useState } from 'react';
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

function App() {
  const [id, setId] = useState<string>('songtomtom');
  const [content, setContent] = useState<string>('Hello~!');

  const [createPost] = useMutation(CREATE_POST, {
    variables: { input: { id } },
  });
  const [createComment] = useMutation(CREATE_COMMENT, {
    variables: { input: { postId: id, content } },
  });
  const { data, loading } = useQuery(LIST_COMMENTS, {
    variables: { where: { postId: id } },
  });

  // console.log('comments: ', data.comments);

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
