import WebSocket from "ws";
import { Message } from "./websocket-server";

// Function to broadcast message to all clients
const broadcastToClients = (
  wss: WebSocket.Server,
  message: Message,
  exclude?: WebSocket
) => {
  wss.clients.forEach((client) => {
    if (client !== exclude && client.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify(message));
    }
  });
};

export default broadcastToClients;
