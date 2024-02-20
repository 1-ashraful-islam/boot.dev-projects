import React, { useState, useEffect } from "react";
import PostCard, { Post } from "./PostCard";

const PostList: React.FC<{
  feed_id: string;
  initialOffset: string;
  initialLimit: string;
}> = ({ feed_id, initialOffset, initialLimit }) => {
  const [fetchError, setFetchError] = useState<string>("");
  // State to manage offset and limit for pagination
  const [offset, setOffset] = useState(initialOffset);
  const limit = initialLimit;
  const [posts, setPosts] = useState<Post[]>([]);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const url = new URL(`http://localhost:8080/v1/posts/${feed_id || ""}`);
        const params = new URLSearchParams({ offset, limit });
        url.search = params.toString();

        const response = await fetch(url.toString());
        if (response.ok) {
          const newPosts: Post[] = await response.json();
          if (newPosts) {
            // Append new posts instead of replacing them
            setPosts((prevPosts) => {
              const uniquePosts = newPosts.filter((newPost) => {
                return !prevPosts.some(
                  (prevPost) => prevPost.id === newPost.id
                );
              });
              return [...prevPosts, ...uniquePosts];
            });
            setFetchError("");
          }
        }
      } catch (error) {
        setFetchError("Error fetching / refreshing Posts");
        console.error("Error fetching Posts", error);
      }
    };

    fetchPosts();
  }, [feed_id, limit, offset]);

  // Load more posts by increasing the offset
  const handleLoadMore = () => {
    setOffset((prevOffset) =>
      (parseInt(prevOffset) + parseInt(limit)).toString()
    );
  };

  return (
    <>
      {fetchError && <div className="network-error">{fetchError}</div>}
      <ul className="horizontal-list">
        {posts.map((post) => (
          <li key={post.id}>
            <PostCard post={post} />
          </li>
        ))}
        <li key="LoadMore">
          <button onClick={handleLoadMore}>Load More</button>
        </li>
      </ul>
    </>
  );
};

export default PostList;
