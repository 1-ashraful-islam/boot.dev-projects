import React, { useState } from "react";
import axios from "axios";

const FeedForm: React.FC = () => {
  const [url, setUrl] = useState("");
  const [apiKey, setApiKey] = useState("");

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const payload = {
      url,
    };

    try {
      const response = await axios.post(
        "http://localhost:8080/v1/feeds",
        payload,
        {
          headers: {
            Authorization: `Bearer ${apiKey}`,
          },
        }
      );

      // Handle response here
      console.log(response.data);
    } catch (error) {
      // Handle error here
      console.error(error);
    }

    // Reset form fields
    setUrl("");
    setApiKey("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        URL:
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
      </label>
      <br />
      <label>
        API Key:
        <input
          type="text"
          value={apiKey}
          onChange={(e) => setApiKey(e.target.value)}
          required
        />
      </label>
      <br />
      <button type="submit">Submit</button>
    </form>
  );
};

export default FeedForm;
