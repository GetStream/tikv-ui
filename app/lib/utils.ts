import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

// Format bytes to human readable compact (KB, MB, GB, TB)
export const formatBytesCompact = (value: number): string => {
  if (value === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(Math.abs(value)) / Math.log(1024));
  const idx = Math.min(i, units.length - 1);
  return `${(value / Math.pow(1024, idx)).toFixed(1)}${units[idx]}`;
};

// Format count to human readable (1K, 10K, 1M, 1B)
export const formatCount = (value: number): string => {
  if (Math.abs(value) < 1000) return value.toFixed(0);
  if (Math.abs(value) < 1000000) return `${(value / 1000).toFixed(1)}K`;
  if (Math.abs(value) < 1000000000) return `${(value / 1000000).toFixed(1)}M`;
  return `${(value / 1000000000).toFixed(1)}B`;
};

// Format time (assumes microseconds input for most TiKV metrics)
export const formatTime = (value: number): string => {
  if (Math.abs(value) < 1000) return `${value.toFixed(0)}µs`;
  if (Math.abs(value) < 1000000) return `${(value / 1000).toFixed(1)}ms`;
  return `${(value / 1000000).toFixed(1)}s`;
};

// Format rate (bytes/s or ops/s)
export const formatRate = (value: number): string => {
  return `${formatBytesCompact(value)}/s`;
};

// Format ratio/percentage
export const formatRatio = (value: number): string => {
  if (value <= 1) return `${(value * 100).toFixed(1)}%`;
  return value.toFixed(2);
};

// Format uptime string (e.g., "2932h5m10s" -> "122d 4h")
// Shows only 2 most significant non-zero units
export const formatUptime = (uptime: string): string => {
  // Extract time components from the uptime string
  const hoursMatch = uptime.match(/(\d+)h/);
  const minsMatch = uptime.match(/(\d+)m(?!o)/); // avoid matching "mo" for months

  const totalHours = hoursMatch ? parseInt(hoursMatch[1], 10) : 0;
  const mins = minsMatch ? parseInt(minsMatch[1], 10) : 0;

  // Convert to total minutes for easier calculation
  const totalMins = totalHours * 60 + mins;

  // Define time units in minutes
  const units = [
    { name: "y", mins: 365 * 24 * 60 },
    { name: "mo", mins: 30 * 24 * 60 },
    { name: "w", mins: 7 * 24 * 60 },
    { name: "d", mins: 24 * 60 },
    { name: "h", mins: 60 },
    { name: "m", mins: 1 },
  ];

  let remaining = totalMins;
  const parts: string[] = [];

  for (const unit of units) {
    if (remaining >= unit.mins) {
      const value = Math.floor(remaining / unit.mins);
      remaining = remaining % unit.mins;
      parts.push(`${value}${unit.name}`);

      // We have the first unit, now find the next non-zero unit
      if (parts.length === 1) {
        for (const nextUnit of units.slice(units.indexOf(unit) + 1)) {
          if (remaining >= nextUnit.mins) {
            const nextValue = Math.floor(remaining / nextUnit.mins);
            parts.push(`${nextValue}${nextUnit.name}`);
            break;
          }
        }
        break;
      }
    }
  }

  return parts.length > 0 ? parts.join(" ") : "0m";
};

// Get formatter based on unit type
export const getFormatter = (unit: string): ((value: number) => string) => {
  switch (unit) {
    case "bytes":
      return formatBytesCompact;
    case "count":
      return formatCount;
    case "time":
      return formatTime;
    case "rate":
      return formatRate;
    case "ratio":
      return formatRatio;
    default:
      return formatCount; // fallback to count formatter
  }
};

// Process gauge data for charts
export const processGaugeData = (points: any[]) => {
  const groupedData: Record<string, any> = {};
  const seriesKeys = new Set<string>();

  points.forEach((point) => {
    const ts = point.ts;
    if (!groupedData[ts]) {
      groupedData[ts] = { ts };
    }

    // Keep raw value - formatting will be done on display
    const value = point.value;

    if (point.labels && Object.keys(point.labels).length > 0) {
      const labelKey = Object.entries(point.labels)
        .map(([k, v]) => `${k}=${v}`)
        .join(", ");

      groupedData[ts][labelKey] = value;
      seriesKeys.add(labelKey);
    } else {
      groupedData[ts]["value"] = value;
      seriesKeys.add("value");
    }
  });

  const data = Object.values(groupedData).sort((a: any, b: any) => a.ts - b.ts);
  return { data, series: Array.from(seriesKeys) };
};

export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export function parseBytes(input: string) {
  if (!input) return 0;

  const match = input
    .trim()
    .match(/^([\d.]+)\s*(B|KB|MB|GB|TB|PB|KiB|MiB|GiB|TiB|PiB)$/i);
  if (!match) throw new Error(`Invalid size: ${input}`);

  const value = parseFloat(match[1]);
  const unit = match[2];

  const multipliers = {
    B: 1,

    KB: 1e3,
    MB: 1e6,
    GB: 1e9,
    TB: 1e12,
    PB: 1e15,

    KiB: 2 ** 10,
    MiB: 2 ** 20,
    GiB: 2 ** 30,
    TiB: 2 ** 40,
    PiB: 2 ** 50,
  };

  return value * multipliers[unit as keyof typeof multipliers];
}
