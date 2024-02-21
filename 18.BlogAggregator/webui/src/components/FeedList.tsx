import React, { useState, useEffect } from "react";
import Modal from "react-modal";
import { FaLink, FaPlus } from "react-icons/fa";
import NewFeedForm from "./NewFeedForm";
import PostList from "./PostList";
import styles from "./FeedList.module.css";
import axios from "axios";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

// Set app element for accessibility; typically set to your root app element
Modal.setAppElement("#root");

interface Feed {
  id: string;
  title: string;
  url: string;
}

const FeedList: React.FC = () => {
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [fetchError, setFetchError] = useState<string>("");
  const [isFormVisible, setIsFormVisible] = useState(false);
  const api_key = "";

  useEffect(() => {
    const fetchFeeds = async () => {
      try {
        const response = await fetch(`http://localhost:8080/v1/feeds`);
        if (response.ok) {
          const data: Feed[] = await response.json();
          if (data) {
            setFeeds(data);
            setFetchError("");
          }
        }
      } catch (error) {
        setFetchError("Error fetching / refreshing feeds");
        console.error("Error fetching feeds", error);
      }
    };

    fetchFeeds();
  }, []);

  const toggleFormVisibility = () => {
    setIsFormVisible(!isFormVisible);
  };

  const handleFeedFollow = async (feed_id: string, api_key: string) => {
    if (!api_key) {
      toast.error("Please login with an API key to follow a feed", {
        position: "top-center",
        autoClose: 3000,
      });
      return;
    }
    try {
      const response = await axios.post(
        `http://localhost:8080/v1/feed_follows`,
        { feed_id },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${api_key}`,
          },
        }
      );
      if (response.status === 200) {
        // Handle response here
        console.log("Feed followed successfully");
        toast.success("Feed followed successfully", {
          position: "top-center",
          autoClose: 3000,
        });
      }
    } catch (error) {
      // Handle error here
      console.error("Error following feed", error);
    }
  };

  return (
    <>
      {fetchError && <div className="network-error">{fetchError}</div>}
      <ul className={styles.feedList}>
        <li>
          <div>
            <h2>
              <button onClick={toggleFormVisibility}>
                <strong>+</strong>{" "}
                {isFormVisible ? "Cancel new feed" : "Add a new Feed"}
              </button>
            </h2>
            {isFormVisible && <NewFeedForm />}
          </div>
        </li>
        {feeds.map((feed) => (
          <li key={feed.id}>
            <h3>
              {feed.title}{" "}
              <a href={feed.url || "#"} target="_blank" rel="noreferrer">
                <FaLink />
              </a>{" "}
              <button onClick={() => handleFeedFollow(feed.id, api_key)}>
                <FaPlus /> Follow the feed
              </button>
            </h3>

            <PostList
              feed_id={feed.id}
              initialOffset={`${0}`}
              initialLimit={`${5}`}
            />
          </li>
        ))}
      </ul>
    </>
  );
};

export default FeedList;
