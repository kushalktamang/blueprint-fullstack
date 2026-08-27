import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./css/index.css";
import App from "./app.tsx";

const rootElement = document.getElementById("root");
if (rootElement === null) {
  throw new Error("root element not found.");
}

createRoot(rootElement).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
