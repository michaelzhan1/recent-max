import { useEffect, useState } from "react";
import type { DataPoint } from "../types/types";
import { Line, LineChart, ReferenceLine, XAxis, YAxis } from "recharts";
import { formatTimestamp } from "../utils/utils";

export default function Data() {
  const [dataArr, setDataArr] = useState<DataPoint[]>([]);
  const [maxValue, setMaxValue] = useState<number | null>(null);

  // data stream
  useEffect(() => {
    const dataEvtSource = new EventSource("http://localhost:8080/stream/data");
    const maxEvtSource = new EventSource("http://localhost:8080/stream/stats");

    dataEvtSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as {
        value: number;
        timestamp: string;
      };

      const now = Date.now();
      setDataArr((prev) => [
        ...prev.filter((point) => now - point.timestamp.getTime() <= 5000),
        {
          value: data.value,
          timestamp: new Date(data.timestamp),
        },
      ]);
    };

    dataEvtSource.onerror = () => {
      dataEvtSource.close();
    };

    maxEvtSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as { maxValue: number | null };
      setMaxValue(data.maxValue);
    };

    maxEvtSource.onerror = () => {
      maxEvtSource.close();
    };

    return () => {
      dataEvtSource.close();
      maxEvtSource.close();
    };
  }, []);

  // make x axis ticks the whole seconds within dataArr
  const firstTime = dataArr.length > 0 ? dataArr[0].timestamp.getTime() : 0;
  const lastTime =
    dataArr.length > 0 ? dataArr[dataArr.length - 1].timestamp.getTime() : 0;
  const xTicks: number[] = [];
  for (let t = Math.floor(lastTime / 1000) * 1000; t > firstTime; t -= 1000) {
    xTicks.push(t);
  }
  xTicks.reverse();

  return (
    <div>
      <h1>Data</h1>
      <div>{maxValue}</div>
      <LineChart
        data={dataArr}
        width={600}
        height={300}
        margin={{ top: 20, right: 80, bottom: 25, left: 20 }}
      >
        <XAxis
          dataKey="timestamp"
          type="number"
          domain={[firstTime, lastTime]}
          tickFormatter={(tick) => formatTimestamp(new Date(tick))}
          interval={0}
          ticks={xTicks}
          label={{
            value: "Timestamp",
            position: "insideBottom",
            offset: -10,
          }}
        />
        <YAxis
          label={{
            value: "Value",
            angle: -90,
            position: "insideLeft",
            offset: 10,
          }}
        />
        <Line
          type="monotone"
          dataKey="value"
          isAnimationActive={false}
          animationDuration={0}
        />
        {maxValue !== null && (
          <ReferenceLine
            y={maxValue}
            stroke="red"
            strokeDasharray="3 3"
            label={{
              value: Math.round(maxValue * 100) / 100,
              position: "right",
              fill: "red",
            }}
          />
        )}
      </LineChart>
    </div>
  );
}
