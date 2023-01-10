import React from 'react';
import './App.css';
import { gql, useSubscription } from '@apollo/client';


const CURRENT_TIME_SUBSCRIPTION = gql`
    subscription OnCurrentTime {
        currentTime {
            unixTime
            timeStamp
        }
    }
`;



function App() {
  const { data, loading } = useSubscription(CURRENT_TIME_SUBSCRIPTION);
  return <h4>New current time: {!loading && data.currentTime.timeStamp}</h4>;
}

export default App;
