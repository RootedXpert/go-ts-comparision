import { randomUUID } from "crypto";
import WebSocket from "ws";
import fs from "fs";
import { Message } from "./websocket-server";
import path from "path";

interface simulateClientsParams {
  NUM_CLIENTS: number;
  MESSAGES_PER_CLIENT: number;
  PORT: number;
}

interface Timing {
  message_id: string;
  time: number;
  client_id: string;
}

interface Matrix {
  client_id: string;
  messages_timings: Timing[];
  avg?: number;
  min?: number;
  max?: number;
}

const simulateClients = ({
  NUM_CLIENTS,
  MESSAGES_PER_CLIENT,
  PORT,
}: simulateClientsParams) => {
  const matric = new Map<string, Matrix>();
  let clientsConnected = 0;
  let clientsCompleted = 0;

  const startSendingMessages = () => {
    matric.forEach((_, client_id) => {
      const client = new WebSocket(`ws://localhost:${PORT}/ws`);

      client.on("open", () => {
        for (let j = 0; j < MESSAGES_PER_CLIENT; j++) {
          const message = {
            sender: client_id,
            id: randomUUID(),
            content: `Message ${j + 1} from ${client_id}`,
            iat: Date.now(),
            type: "message",
          };

          client.send(JSON.stringify(message));
        }
      });

      client.on("message", (message: WebSocket.Data) => {
        try {
          const parsedMessage = JSON.parse(message.toString()) as Message;
          const timeTaken = Date.now() - parsedMessage.iat;

          const existingMatrix = matric.get(client_id)!;
          existingMatrix.messages_timings.push({
            message_id: parsedMessage.id,
            time: timeTaken,
            client_id: parsedMessage.sender,
          });

          // Check if all messages have been received
          if (existingMatrix.messages_timings.length === MESSAGES_PER_CLIENT) {
            client.close();
          }
        } catch (error) {
          console.error(`Error processing message from ${client_id}: ${error}`);
        }
      });

      client.on("close", () => {
        clientsCompleted++;

        if (clientsCompleted === NUM_CLIENTS) {
          // All clients have completed sending messages
          try {
            calculateOverallAverages();
          } catch (error) {}
        }
      });

      client.on("error", (error) => {
        console.error(`Client ${client_id}: Error - ${error.message}`);
      });
    });
  };

  const calculateOverallAverages = () => {
    const results: { [client_id: string]: Matrix } = {};

    matric.forEach((value, client_id) => {
      const timings = value.messages_timings;
      let sum = 0;
      let min = Number.MAX_VALUE;
      let max = Number.MIN_VALUE;

      timings.forEach((timing) => {
        sum += timing.time;
        if (timing.time < min) {
          min = timing.time;
        }
        if (timing.time > max) {
          max = timing.time;
        }
      });

      const avg = sum / timings.length;

      // Update matrix with averages, min, max
      const updatedMatrix = {
        ...value,
        avg,
        min,
        max,
      };

      results[client_id] = updatedMatrix;
    });

    // Convert results object to JSON string
    const jsonString = JSON.stringify(results, null, 2);

    // Create metric directory if it doesn't exist
    const metricDir = path.join(process.cwd(), "metric");
    if (!fs.existsSync(metricDir)) {
      fs.mkdirSync(metricDir);
    }

    // Construct file path with timestamp
    const filePath = path.join(
      metricDir,

      `results-client-typescript-${NUM_CLIENTS}-messages-${MESSAGES_PER_CLIENT}.json`
    );

    fs.writeFile(filePath, jsonString, (err) => {
      if (err) {
        console.error(`Error writing to file ${filePath}: ${err}`);
        process.exit(1);
      } else {
        console.log(`Results saved to ${filePath}`);
        process.exit(0);
      }
    });
  };

  for (let i = 0; i < NUM_CLIENTS; i++) {
    const client_id = `Client_${i + 1}`;

    // Initialize client matrix entry
    matric.set(client_id, {
      client_id,
      messages_timings: [],
    });

    const client = new WebSocket(`ws://localhost:${PORT}/ws`);

    client.on("open", () => {
      clientsConnected++;

      if (clientsConnected === NUM_CLIENTS) {
        // All clients are now connected, start sending messages
        startSendingMessages();
      }
    });

    client.on("error", (error) => {
      console.error(`Client ${client_id}: Error - ${error.message}`);
    });
  }
};

export default simulateClients;
