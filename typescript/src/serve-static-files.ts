import * as http from "http";
import * as path from "path";
import * as fs from "fs/promises";

// Function to serve static files
const serveStaticFile = async (
  req: http.IncomingMessage,
  res: http.ServerResponse,
  HTML_FILE_PATH: string
) => {
  try {
    let filePath = path.join(__dirname, req.url as string);
    if (req.url === "/" || req.url === "/index.html") {
      filePath = HTML_FILE_PATH;
    }
    const data = await fs.readFile(filePath);
    res.writeHead(200);
    res.end(data);
  } catch (err) {
    res.writeHead(404, { "Content-Type": "text/plain" });
    res.end("File not found");
  }
};

export default serveStaticFile;
