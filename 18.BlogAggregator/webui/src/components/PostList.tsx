import React, { useState, useEffect } from "react";

interface Post {
  id: string;
  title: string;
  url: string;
}

const PostList: React.FC<{
  feed_id: string;
  offset: string;
  limit: string;
}> = ({ feed_id, offset, limit }) => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [fetchError, setFetchError] = useState<string>("");

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const url = new URL(`http://localhost:8080/v1/posts/${feed_id || ""}`);
        const params = new URLSearchParams();

        if (offset) {
          params.append("offset", offset);
        }

        if (limit) {
          params.append("limit", limit);
        }

        url.search = params.toString();

        const response = await fetch(url.toString());
        if (response.ok) {
          const data: Post[] = await response.json();
          if (data) {
            setPosts(data);
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

  return (
    <>
      <div>
        <h2>Posts</h2>
      </div>
      {fetchError && <div className="network-error">{fetchError}</div>}
      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            {post.title} ({post.url})
          </li>
        ))}
      </ul>
    </>
  );
};

export default PostList;
