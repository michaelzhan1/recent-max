import { useMemo, useState } from "react";
import { Controls } from "./components/controls";
import Data from "./components/data";
import type { DataPoint, Stat } from "./types/types";
import { MonotonicDeque } from "./utils/deque";
import { DequeContext } from "./context/app-context";

export default function App() {
  const maxDeque = useMemo(() => new MonotonicDeque("max"), []);
  const minDeque = useMemo(() => new MonotonicDeque("min"), []);

  const [dataArr, setDataArr] = useState<DataPoint[]>([]);
  const [stats, setStats] = useState<Stat>({
    maxValue: null,
    minValue: null,
    avg: null,
  });

  return (
    <DequeContext.Provider value={{ maxDeque, minDeque }}>
      <div>
        <Data
          dataArr={dataArr}
          setDataArr={setDataArr}
          stats={stats}
          setStats={setStats}
        />
        <Controls setDataArr={setDataArr} setStats={setStats} />
      </div>
    </DequeContext.Provider>
  );
}
