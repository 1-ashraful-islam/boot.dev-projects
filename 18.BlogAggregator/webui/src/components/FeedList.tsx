import React, { useState, useEffect } from 'react';

interface Feed {
  id: number;
  title: string;
  // Add other feed properties as needed
}

const FeedList: React.FC = () => {
  const [feeds, setFeeds] = useState<Feed[]>([]);

  useEffect(() => {
    const fetchFeeds = async () => {
      // Replace with your actual feed API call
      const response = await fetch('http://localhost:8080/feeds');
      if (response.ok) {
        const data: Feed[] = await response.json();
        setFeeds(data);
      }
    };

    fetchFeeds();
  }, []);

  return (
    <ul>
      {feeds.map((feed) => (
        <li key={feed.id}>{feed.title}</li>
      ))}
    </ul>
  );
};

export default FeedList;
