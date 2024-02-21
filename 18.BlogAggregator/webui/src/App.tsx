import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import "./App.css";
import LandingPage from "./pages/LandingPage";
import APIListPage from "./pages/APIListPage";
import { ToastContainer } from "react-toastify";
import ExplorePage from "./pages/ExplorePage";

function Navigation() {
  return (
    <nav>
      <ul className="App-nav">
        <div className="nav-links">
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/discover">Discover</Link>
          </li>
          <li>
            <Link to="/api">API Reference</Link>
          </li>
        </div>
        <div className="auth-links">
          <li>
            <Link to="/login">Login</Link>
          </li>
          <li>
            <Link to="/signup">Sign Up</Link>
          </li>
        </div>
      </ul>
    </nav>
  );
}

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <Navigation />
        </header>
        <ToastContainer />
        <main className="App-main">
          <Routes>
            <Route path="/" element={<LandingPage />} />
            <Route path="/discover" element={<ExplorePage />} />
            <Route path="/api" element={<APIListPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
