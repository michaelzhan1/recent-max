import { useState } from "react";

export function Controls() {
  const [paused, setPaused] = useState<boolean>(false);

  const handlePauseToggle = () => {
    if (paused) {
      fetch("http://localhost:8080/resume", { method: "POST" })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to resume data generation");
          }
          setPaused(false);
        })
        .catch((error) => {
          console.error("Error resuming data generation:", error);
        });
    } else {
      fetch("http://localhost:8080/pause", { method: "POST" })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to pause data generation");
          }
          setPaused(true);
        })
        .catch((error) => {
          console.error("Error pausing data generation:", error);
        });
    }
  };

  return (
    <div style={{ display: "flex", marginTop: "20px" }}>
      <button onClick={handlePauseToggle} style={{ padding: "10px 20px", fontSize: "16px" }}>
        {paused ? "Resume" : "Pause"}
      </button>
    </div>
  )
}
