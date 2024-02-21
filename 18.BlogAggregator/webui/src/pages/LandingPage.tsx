import FeedList from "../components/FeedList";
import React, { useState } from "react";
import NotLoggedInLanding from "../components/NotLoggedInLanding";

const LandingPage: React.FC = () => {
  const [loggedInUser, setLoggedInUser] = useState(false);

  return <>{loggedInUser ? <FeedList /> : <NotLoggedInLanding />}</>;
};

export default LandingPage;
