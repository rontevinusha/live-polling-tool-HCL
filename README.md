\# Live Polling Tool



A real-time polling application built using React, Go, MongoDB, Redis, and WebSocket.



\## Features



\- Admin login with JWT authentication

\- Create polls with multiple options

\- Vote on polls

\- Real-time vote updates using WebSocket

\- Vote percentage display

\- MongoDB for poll data storage

\- Redis for fast vote counting

\- REST API using Go and Gin

\- Responsive React frontend



\## Tech Stack



\### Frontend

\- React

\- Vite

\- JavaScript

\- CSS



\### Backend

\- Go

\- Gin

\- JWT Authentication

\- Gorilla WebSocket



\### Database \& Cache

\- MongoDB

\- Redis / Memurai



\## Project Structure



```text

live-polling-tool/

├── backend/

│   ├── controllers/

│   ├── middleware/

│   ├── models/

│   ├── services/

│   ├── websocket/

│   ├── main.go

│   ├── go.mod

│   └── go.sum

│

├── frontend/

│   ├── src/

│   ├── package.json

│   └── vite.config.js

│

└── README.md

