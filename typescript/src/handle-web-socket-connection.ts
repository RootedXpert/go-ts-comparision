import WebSocket from "ws";
import broadcastToClients from "./broad-cast-client";
import { Message } from "./websocket-server";
import { randomUUID } from "crypto";

// Function to handle WebSocket connections
const handleWebSocketConnection = (wss: WebSocket.Server) => {
  wss.on("connection", (ws: WebSocket) => {
    broadcastToClients(wss, {
      sender: "Server",
      content: "A new user has joined",
      id: randomUUID(),
      iat: Date.now(),
      type: "join",
    });

    ws.on("message", (message: Buffer) => {
      const parsedMessage: Message = JSON.parse(message.toString());
      broadcastToClients(wss, parsedMessage, ws);
    });

    ws.on("close", () => {
      broadcastToClients(wss, {
        sender: "Server",
        content: "A user has left",
        id: randomUUID(),
        iat: Date.now(),
        type: "left",
      });
    });
  });
};

export default handleWebSocketConnection;
