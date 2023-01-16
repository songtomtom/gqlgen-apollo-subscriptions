import React from 'react';
import './App.css';
import { gql, useMutation, useQuery } from '@apollo/client';

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

const POST_ID = 'fe440985-3dad-4f47-968b-41668ca7f03c';

function App() {
  const [createPost, { data: post }] = useMutation(CREATE_POST, {
    variables: { postId: POST_ID },
  });
  const [createComment, { data: comment }] = useMutation(CREATE_COMMENT, {
    variables: { postId: POST_ID },
  });

  const { data } = useQuery(LIST_COMMENTS, {
    variables: { postId: POST_ID },
  });

  console.log('comments: ', data);

  const onCreatePost = async () => {
    const { data } = await createPost();
    if (data) {
      console.log('creat post: ', data);
    }
  };

  const onCreateComments = async () => {
    const { data } = await createComment();
    if (data) {
      console.log('creat comment: ', data);
    }
  };

  return (
    <div>
      <button onClick={onCreatePost}>create post</button>
      <button onClick={onCreateComments}>create comment</button>
      <li>
        {data &&
          data.map((item: any) => {
            console.log(item);
            return <ul>{item.id}</ul>;
          })}
      </li>
    </div>
  );
}

export default App;
