import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import "./App.css";
import LandingPage from "./pages/LandingPage";
import APIListPage from "./pages/APIListPage";

function Navigation() {
  return (
    <header className="App-header">
      <nav>
        <ul className="App-nav">
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/apilist">API Reference</Link>
          </li>
        </ul>
      </nav>
    </header>
  );
}

function App() {
  return (
    <Router>
      <div className="App">
        <Navigation />
        <main className="App-main">
          <Routes>
            <Route path="/apilist" element={<APIListPage />} />
            <Route path="/" element={<LandingPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
