import { FaLink } from "react-icons/fa";
import styles from "./PostCard.module.css";

export interface Post {
  id: string;
  title: string;
  url: string;
  description: string;
}

export default function PostCard({ post }: { post: Post }) {
  return (
    <div className={styles.postCardContainer}>
      <h4>
        {post.title || "Post Title"}{" "}
        <a href={post.url || "#"} target="_blank" rel="noreferrer">
          <FaLink />
        </a>
      </h4>
      <p>
        {post.description && post.description.length > 300
          ? post.description.slice(0, 297) + "..."
          : post.description || "Post Description"}
      </p>
    </div>
  );
}
