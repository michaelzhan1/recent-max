export interface Stat {
  maxValue: number | null;
  minValue: number | null;
  avg: number | null;
}

export interface DataPoint {
  value: number;
  timestamp: Date;
}
