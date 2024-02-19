import React, { useState, useEffect } from "react";
import Modal from "react-modal";
import NewFeedForm from "./NewFeedForm";

// Set app element for accessibility; typically set to your root app element
Modal.setAppElement("#root");

interface Feed {
  id: number;
  title: string;
  url: string;
}

const modalStyles = {
  content: {
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
  },
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.75)",
  },
};

const FeedList: React.FC = () => {
  const [feeds, setFeeds] = useState<Feed[]>([]);
  const [modalIsOpen, setModalIsOpen] = useState(false);

  const [isFormVisible, setIsFormVisible] = useState(false);

  useEffect(() => {
    const fetchFeeds = async () => {
      const response = await fetch(`http://localhost:8080/v1/feeds`);
      if (response.ok) {
        const data: Feed[] = await response.json();
        setFeeds(data);
      }
    };

    fetchFeeds();
  }, []);

  const handleOpenModal = () => {
    setModalIsOpen(true);
  };

  const handleCloseModal = () => {
    setModalIsOpen(false);
  };

  const toggleFormVisibility = () => {
    setIsFormVisible(!isFormVisible);
  };

  return (
    <>
      <Modal
        isOpen={modalIsOpen}
        onRequestClose={handleCloseModal}
        contentLabel="Add New Feed"
        style={modalStyles}
      >
        <NewFeedForm />
        <button onClick={handleCloseModal}>Close</button>
      </Modal>
      <div>
        <h2>
          Feeds{" "}
          <button onClick={toggleFormVisibility}>
            <strong>+</strong>{" "}
            {isFormVisible ? "Cancel new feed" : "Add a new Feed"}
          </button>
          <button onClick={handleOpenModal}>
            <strong>+</strong>Add a new Feed with Modal
          </button>
        </h2>
        {isFormVisible && <NewFeedForm />}
      </div>
      <ul>
        {feeds.map((feed) => (
          <li key={feed.id}>
            {feed.title} ({feed.url})
          </li>
        ))}
      </ul>
    </>
  );
};

export default FeedList;
