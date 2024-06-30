import { Server } from "http";

import { Server as WebSocket } from "ws";

// Function to handle WebSocket upgrade
const handleWebSocketUpgrade = (
  server: Server,
  wss: WebSocket,
  WS_PATH: string
) => {
  server.on("upgrade", (request, socket, head) => {
    const url = new URL(`http://${request.headers.host}${request.url}`);
    if (url.pathname === WS_PATH) {
      wss.handleUpgrade(request, socket, head, (ws) => {
        wss.emit("connection", ws, request);
      });
    } else {
      socket.destroy();
    }
  });
};

export default handleWebSocketUpgrade;
