import { useEffect, useState } from "react";

export default function Data() {
  const [data, setData] = useState<number[]>([]);

  useEffect(() => {
    const evtSource = new EventSource("http://localhost:8080/stream/data");

    evtSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as { value: number };
      setData((prev) => [...prev, data.value]);
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
      <h1>Data</h1>
      <p>{data.join(", ")}</p>
    </div>
  );
}
