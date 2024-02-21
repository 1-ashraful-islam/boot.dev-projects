import FeedList from "./FeedList";

export default function NotLoggedInLanding() {
  return (
    <>
      <h1> Welcome to blog aggregator!</h1>
      <p>
        {" "}
        Keep up with your favorite blogs, find posts that resonates with you or
        explore posts from a curated list of rss feeds
      </p>
      <p>
        Log in to see posts from your feeds you follow, or sign up if you are
        new here!
      </p>
      <h1> Discover New RSS feeds</h1>
      <FeedList />
    </>
  );
}
