import { useState } from "react";
import { Controls } from "./components/controls";
import Data from "./components/data";
import type { DataPoint, Stat } from "./types/types";

export default function App() {
  const [dataArr, setDataArr] = useState<DataPoint[]>([]);
  const [stats, setStats] = useState<Stat>({
    maxValue: null,
    minValue: null,
    avg: null,
  });

  return (
    <div>
      <Data
        dataArr={dataArr}
        setDataArr={setDataArr}
        stats={stats}
        setStats={setStats}
      />
      <Controls setDataArr={setDataArr} setStats={setStats} />
    </div>
  );
}
