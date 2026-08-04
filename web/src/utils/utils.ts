export function formatTimestamp(timestamp: Date): string {
  const hours = timestamp.getHours().toString().padStart(2, "0");
  const minutes = timestamp.getMinutes().toString().padStart(2, "0");
  const seconds = timestamp.getSeconds().toString().padStart(2, "0");
  return `${hours}:${minutes}:${seconds}`;
}