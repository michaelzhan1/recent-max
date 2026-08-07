import { useContext, useEffect } from "react";
import { Line, LineChart, ReferenceLine, XAxis, YAxis } from "recharts";
import { DequeContext } from "../context/app-context";
import type { DataPoint, Stat } from "../types/types";
import { formatTimestamp } from "../utils/utils";

interface DataProps {
  dataArr: DataPoint[];
  setDataArr: React.Dispatch<React.SetStateAction<DataPoint[]>>;
  stats: Stat;
  setStats: React.Dispatch<React.SetStateAction<Stat>>;
}

export default function Data({ dataArr, setDataArr, stats, setStats }: DataProps) {
  const { maxDeque, minDeque } = useContext(DequeContext) ?? {};

  // data stream
  useEffect(() => {
    const dataEvtSource = new EventSource("http://localhost:8080/stream/data");

    dataEvtSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as {
        value: number;
        timestamp: string;
      };
      const dataPoint: DataPoint = {
        value: data.value,
        timestamp: new Date(data.timestamp),
      };

      maxDeque?.push(dataPoint);
      minDeque?.push(dataPoint);

      const currentTime = dataPoint.timestamp.getTime();
      setDataArr((prev) => {
        let i = 0;
        while (i < prev.length && currentTime - prev[i].timestamp.getTime() > 5000) {
          i++;
        }
        const newArr = [...prev.slice(i), dataPoint];

        const sum = newArr.reduce((acc, dp) => acc + dp.value, 0);
        const avg = newArr.length > 0 ? sum / newArr.length : null;

        setStats({
          maxValue: maxDeque?.peek()?.value ?? null,
          minValue: minDeque?.peek()?.value ?? null,
          avg: avg,
        });

        return newArr;
      });
    };

    dataEvtSource.onerror = () => {
      dataEvtSource.close();
    };

    return () => {
      dataEvtSource.close();
    };
  }, [setDataArr, setStats, maxDeque, minDeque]);

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
        {stats.maxValue !== null && (
          <ReferenceLine
            y={stats.maxValue}
            stroke="red"
            strokeDasharray="3 3"
            label={{
              value: `Max: ${Math.round(stats.maxValue * 100) / 100}`,
              position: "right",
              fill: "red",
            }}
          />
        )}
        {stats.minValue !== null && (
          <ReferenceLine
            y={stats.minValue}
            stroke="blue"
            strokeDasharray="3 3"
            label={{
              value: `Min: ${Math.round(stats.minValue * 100) / 100}`,
              position: "right",
              fill: "blue",
            }}
          />
        )}
        {stats.avg !== null && (
          <ReferenceLine
            y={stats.avg}
            stroke="green"
            strokeDasharray="3 3"
            label={{
              value: `Avg: ${Math.round(stats.avg * 100) / 100}`,
              position: "right",
              fill: "green",
            }}
          />
        )}
      </LineChart>
    </div>
  );
}
