import { createContext } from "react";
import type { MonotonicDeque } from "../utils/deque";

interface DequeContextType {
  maxDeque: MonotonicDeque;
  minDeque: MonotonicDeque;
}

export const DequeContext = createContext<DequeContextType | null>(null);