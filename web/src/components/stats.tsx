import { useEffect, useState } from "react";

export default function Stats() {
  const [maxValue, setMaxValue] = useState<number | null>(null);

  useEffect(() => {
    const evtSource = new EventSource("http://localhost:8080/stream/stats");

    evtSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as { maxValue: number };
      setMaxValue(data.maxValue);
    };

    evtSource.onerror = () => {
      evtSource.close();
    };

    return () => {
      evtSource.close();
    };
  }, []);

  return (
    <div>
      <h1>Max Value</h1>
      <p>{maxValue !== null ? maxValue : "No values yet"}</p>
    </div>
  );
}
