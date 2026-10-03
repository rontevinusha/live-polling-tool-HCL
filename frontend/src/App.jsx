import { useEffect, useState } from "react";
import "./App.css";

function App() {
  const [polls, setPolls] = useState([]);

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState(localStorage.getItem("token") || "");

  const [question, setQuestion] = useState("");
  const [option1, setOption1] = useState("");
  const [option2, setOption2] = useState("");

  const [message, setMessage] = useState("");

  const loadPolls = async () => {
    try {
      const response = await fetch("http://localhost:8080/polls");
      const data = await response.json();
      setPolls(data);
    } catch (error) {
      console.error("Failed to load polls:", error);
    }
  };

  const login = async () => {
    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        setMessage(data.error || "Login failed");
        return;
      }

      localStorage.setItem("token", data.token);
      setToken(data.token);
      setMessage("Login successful!");
      setPassword("");
    } catch (error) {
      setMessage("Login failed");
    }
  };

  const createPoll = async () => {
    if (!token) {
      setMessage("Please login first");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/polls", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          question,
          options: [
            {
              id: "1",
              text: option1,
              votes: 0,
            },
            {
              id: "2",
              text: option2,
              votes: 0,
            },
          ],
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        setMessage(data.error || "Failed to create poll");
        return;
      }

      setMessage("Poll created successfully!");
      setQuestion("");
      setOption1("");
      setOption2("");

      loadPolls();
    } catch (error) {
      setMessage("Failed to create poll");
    }
  };

  const vote = async (pollId, optionId) => {
    try {
      await fetch(
        `http://localhost:8080/polls/${pollId}/vote/${optionId}`,
        {
          method: "POST",
        }
      );

      loadPolls();
    } catch (error) {
      console.error("Vote failed:", error);
    }
  };

  useEffect(() => {
    loadPolls();

    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.onmessage = () => {
      loadPolls();
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div className="app">
      <header>
        <h1>Live Polling Tool</h1>
        <p>Vote and see poll results in real time</p>
      </header>

      <main>
        {!token ? (
          <div className="poll-card">
            <h2>Admin Login</h2>

            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />

            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />

            <button className="create-button" onClick={login}>
              Login
            </button>
          </div>
        ) : (
          <div className="poll-card">
            <h2>Create a Poll</h2>

            <input
              type="text"
              placeholder="Poll question"
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
            />

            <input
              type="text"
              placeholder="Option 1"
              value={option1}
              onChange={(e) => setOption1(e.target.value)}
            />

            <input
              type="text"
              placeholder="Option 2"
              value={option2}
              onChange={(e) => setOption2(e.target.value)}
            />

            <button className="create-button" onClick={createPoll}>
              Create Poll
            </button>
          </div>
        )}

        {message && <p className="message">{message}</p>}

        {polls.length === 0 ? (
          <p className="empty">No polls available.</p>
        ) : (
          polls.map((poll) => (
            <div className="poll-card" key={poll.id}>
              <h2>{poll.question}</h2>

              {poll.options.map((option) => (
                <button
                  className="option-button"
                  key={option.id}
                  onClick={() => vote(poll.id, option.id)}
                >
                  <span>{option.text}</span>
                  <strong>{option.votes} votes</strong>
                </button>
              ))}
            </div>
          ))
        )}
      </main>
    </div>
  );
}

export default App;
