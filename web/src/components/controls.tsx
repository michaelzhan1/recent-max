import { useState } from "react";
import type { DataPoint, Stat } from "../types/types";

interface ControlsProps {
  setDataArr: React.Dispatch<React.SetStateAction<DataPoint[]>>;
  setStats: React.Dispatch<React.SetStateAction<Stat>>;
}

async function pause() {
  return fetch("http://localhost:8080/pause", { method: "POST" });
}

async function resume() {
  return fetch("http://localhost:8080/resume", { method: "POST" });
}

export function Controls({ setDataArr, setStats }: ControlsProps) {
  const [paused, setPaused] = useState<boolean>(false);

  const handleTogglePause = () => {
    if (paused) {
      resume()
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
      pause()
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

  const handleClear = () => {
    setDataArr([]);
    setStats({
      maxValue: null,
      minValue: null,
      avg: null,
    });
  };

  const handleFormSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const value = Number(formData.get("value"));
    const mu = Number(formData.get("mu"));
    const sigma = Number(formData.get("sigma"));

    if (Number.isNaN(value) || Number.isNaN(mu) || Number.isNaN(sigma)) {
      alert("Please enter valid numbers for value, mu, and sigma.");
      return;
    }

    fetch("http://localhost:8080/reset", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ value, mu, sigma }),
    })
      .then((response) => {
        if (!response.ok) {
          alert("Failed to reset data generation");
          console.error("Failed to reset data generation");
        } else {
          resume()
            .then((response) => {
              if (!response.ok) {
                alert("Failed to resume data generation");
                console.error("Failed to resume data generation");
              }
              handleClear();
              setPaused(false);
            })
            .catch((error) => {
              alert("Error resuming data generation:");
              console.error("Error resuming data generation:", error);
            });
        }
      })
      .catch((error) => {
        alert("Error resetting data generation:");
        console.error("Error resetting data generation:", error);
      });
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "flex-start",
        gap: "10px",
      }}
    >
      <div style={{ display: "flex", marginTop: "20px" }}>
        <button onClick={handleTogglePause}>
          {paused ? "Resume" : "Pause"}
        </button>
        <button onClick={handleClear}>Clear</button>
      </div>
      <form
        onSubmit={handleFormSubmit}
        style={{ display: "flex", flexDirection: "column", gap: "6px" }}
      >
        <label htmlFor="value">Value</label>
        <input
          type="number"
          name="value"
          placeholder="Value"
          step="any"
          required
        />
        <label htmlFor="mu">Mu</label>
        <input type="number" name="mu" placeholder="Mu" step="any" required />
        <label htmlFor="sigma">Sigma</label>
        <input
          type="number"
          name="sigma"
          placeholder="Sigma"
          step="any"
          required
        />
        <button type="submit">Reset</button>
      </form>
    </div>
  );
}
