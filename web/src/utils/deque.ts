import type { DataPoint } from "../types/types";

export class MonotonicDeque {
  arr: DataPoint[];
  dir: "max" | "min";

  constructor(dir: "max" | "min" = "max") {
    this.arr = [];
    this.dir = dir;
  }

  push(x: DataPoint): void {
    const comp =
      this.dir === "max"
        ? (a: number, b: number) => a < b
        : (a: number, b: number) => a > b;
    while (
      this.arr.length > 0 &&
      comp(this.arr[this.arr.length - 1].value, x.value)
    ) {
      this.arr.pop();
    }
    this.arr.push(x);
  }

  peek(): DataPoint | undefined {
    if (
      this.arr.length > 0 &&
      Date.now() - this.arr[0].timestamp.getTime() > 5000
    ) {
      this.arr.shift();
    }

    if (this.arr.length > 0) {
      return this.arr[0];
    } else {
      return undefined;
    }
  }
}
