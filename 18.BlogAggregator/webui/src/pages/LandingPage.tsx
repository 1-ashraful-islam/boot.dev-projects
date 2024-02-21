import { useAuth } from "../components/AuthContext";
import FeedList from "../components/FeedList";
import NotLoggedInLanding from "../components/NotLoggedInLanding";

const LandingPage: React.FC = () => {
  const { isLoggedIn } = useAuth();
  return <>{isLoggedIn ? <FeedList /> : <NotLoggedInLanding />}</>;
};

export default LandingPage;
