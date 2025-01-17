import * as http from "http";

import * as path from "path";
import WebSocket from "ws";
import serveStaticFile from "./serve-static-files";
import handleWebSocketUpgrade from "./handle-Web-Socket-Upgrade";
import handleWebSocketConnection from "./handle-web-socket-connection";

// Constants
const PORT = 8080;
const HTML_FILE_PATH = path.join(process.cwd(), "static", "index.html");
const WS_PATH = "/ws";

export interface Message {
  sender: string;
  content: string;
  id: string;
  iat: number;
  type: string;
}

// Create HTTP server
const server = http.createServer((req, res) => {
  if (req.method === "GET") {
    serveStaticFile(req, res, HTML_FILE_PATH).catch((err) => {
      console.error("Error serving static file:", err);
      res.writeHead(500, { "Content-Type": "text/plain" });
      res.end("Internal Server Error");
    });
  }
});

// Create WebSocket server
const wss = new WebSocket.Server({ noServer: true });

// Handle WebSocket upgrade requests
handleWebSocketUpgrade(server, wss, WS_PATH);

// Handle WebSocket connections
handleWebSocketConnection(wss);

// Start HTTP server
server.listen(PORT, () => {
  console.log(`Server started on http://localhost:${PORT}`);
});
