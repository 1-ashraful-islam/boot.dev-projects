import FeedList from "./FeedList";

export default function NotLoggedInLanding() {
  return (
    <>
      <div
        style={{
          // display: "flex",
          justifyContent: "center",
          textAlign: "center",
          fontSize: "1.5em",
          maxWidth: "1000px",
          margin: "auto",
        }}
      >
        <h1 style={{ fontSize: "2.5em" }}>Welcome to the blog aggregator!</h1>
        <p>
          Keep up with your favorite blogs, find new blogs that resonate with
          you or explore posts from a curated list of rss feeds
        </p>
        <p>
          Log in to see posts from your curated list, or sign up if you are new
          here!
        </p>
        <br />
        <h2> Discover New RSS feeds</h2>
      </div>
      <FeedList />
    </>
  );
}
