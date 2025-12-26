"use client";

import { useState, useMemo } from "react";
import { Area, AreaChart, Line, LineChart, XAxis, YAxis } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { MetricCard } from "@/components/ui/metric-card";
import { InfoIcon } from "lucide-react";
import dayjs from "dayjs";
import { getFormatter, processGaugeData } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface GaugePoint {
  ts: number;
  value: number;
  labels: Record<string, string>;
}

interface Gauge {
  name: string;
  description: string;
  unit: string;
  points: GaugePoint[];
  labels?: Record<string, string>[];
}

interface Store {
  store: {
    id: number;
    status_address: string;
  };
}

interface GaugeChartProps {
  gauge: Gauge;
  index: number;
  stores: Store[];
}

export function GaugeChart({ gauge, index, stores }: GaugeChartProps) {
  // Extract dynamic label keys (excluding 'instance' which comes from PD stores)
  const labelOptions = useMemo(() => {
    const labelsSource = gauge.labels || gauge.points.map((p) => p.labels);
    const labelMap: Record<string, Set<string>> = {};

    for (const labels of labelsSource) {
      if (!labels) continue;
      for (const [key, value] of Object.entries(labels)) {
        // Skip 'instance' - it uses PD stores data
        if (key === "instance") continue;
        if (!labelMap[key]) {
          labelMap[key] = new Set();
        }
        labelMap[key].add(value);
      }
    }

    // Convert to array format and sort keys
    return Object.entries(labelMap)
      .map(([key, values]) => ({
        key,
        values: Array.from(values).sort(),
      }))
      .sort((a, b) => a.key.localeCompare(b.key));
  }, [gauge.labels, gauge.points]);

  // Instance values from PD stores
  const instanceValues = useMemo(
    () => stores.map((s) => s.store.status_address),
    [stores]
  );

  // State for all filters - instance + dynamic labels
  const [filters, setFilters] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = { instance: "all" };
    labelOptions.forEach(({ key }) => {
      initial[key] = "all";
    });
    return initial;
  });

  const updateFilter = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  // Filter points by all active filters
  const filteredPoints = useMemo(() => {
    return gauge.points.filter((point) => {
      for (const [key, filterValue] of Object.entries(filters)) {
        if (filterValue !== "all" && point.labels?.[key] !== filterValue) {
          return false;
        }
      }
      return true;
    });
  }, [gauge.points, filters]);

  const { data: chartData, series } = processGaugeData(filteredPoints);
  const formatter = getFormatter(gauge.unit);

  const dynamicConfig: ChartConfig = {
    ...series.reduce((acc, key, idx) => {
      const colorVar = `--chart-${(idx % 5) + 1}`;
      acc[key] = {
        label: key,
        color: `var(${colorVar})`,
      };
      return acc;
    }, {} as Record<string, { label: string; color: string }>),
  };

  return (
    <MetricCard>
      <div className="font-semibold text-gray-500 mb-2 flex items-center justify-between">
        <span>{gauge.name}</span>
      </div>
      {(instanceValues.length > 0 || labelOptions.length > 0) && (
        <div className="flex flex-wrap gap-2 mb-10">
          {instanceValues.length > 0 && (
            <Select
              value={filters.instance || "all"}
              onValueChange={(value) => updateFilter("instance", value)}
            >
              <SelectTrigger
                className="h-7 w-auto rounded-lg pl-2.5 cursor-pointer text-xs"
                aria-label="Filter by instance"
                data-size="xs"
              >
                <span className="text-muted-foreground mr-1">instance:</span>
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent align="end" className="rounded-xl max-h-60">
                <SelectItem
                  value="all"
                  className="rounded-lg [&_span]:flex cursor-pointer"
                >
                  <div className="flex items-center gap-2 text-xs">All</div>
                </SelectItem>
                {instanceValues.map((addr) => (
                  <SelectItem
                    key={addr}
                    value={addr}
                    className="rounded-lg [&_span]:flex cursor-pointer"
                  >
                    <div className="flex items-center gap-2 text-xs">
                      {addr}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          {labelOptions.map(({ key, values }) => (
            <Select
              key={key}
              value={filters[key] || "all"}
              onValueChange={(value) => updateFilter(key, value)}
            >
              <SelectTrigger
                className="h-7 w-auto rounded-lg pl-2.5 cursor-pointer text-xs"
                aria-label={`Filter by ${key}`}
                data-size="xs"
              >
                <span className="text-muted-foreground mr-1">{key}:</span>
                <SelectValue placeholder="All" />
              </SelectTrigger>
              <SelectContent align="end" className="rounded-xl max-h-60">
                <SelectItem
                  value="all"
                  className="rounded-lg [&_span]:flex cursor-pointer"
                >
                  <div className="flex items-center gap-2 text-xs">All</div>
                </SelectItem>
                {values.map((value) => (
                  <SelectItem
                    key={value}
                    value={value}
                    className="rounded-lg [&_span]:flex cursor-pointer"
                  >
                    <div className="flex items-center gap-2 text-xs">
                      {value}
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          ))}
        </div>
      )}
      <ChartContainer config={dynamicConfig}>
        {series.length > 1 ? (
          <LineChart data={chartData} margin={{ left: -10 }}>
            <ChartTooltip
              content={
                <ChartTooltipContent
                  formatter={(value, name, item) => (
                    <div className="flex flex-1 justify-between items-center gap-4">
                      <div
                        className="flex items-center gap-2"
                        title={name as string}
                      >
                        <div
                          className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                          style={{ backgroundColor: item.color }}
                        />
                        <span className="text-muted-foreground truncate max-w-[300px]">
                          {name}
                        </span>
                      </div>
                      <span className="font-mono font-medium">
                        {formatter(value as number)}
                      </span>
                    </div>
                  )}
                />
              }
            />
            <XAxis
              dataKey="ts"
              tickFormatter={(value) => dayjs(value).format("HH:mm")}
            />
            <YAxis tickFormatter={formatter} width={70} tickMargin={2} />
            {series.slice(0, 10).map((key) => (
              <Line
                key={key}
                dataKey={key}
                type="monotone"
                stroke={dynamicConfig[key]?.color}
                strokeWidth={2}
                dot={false}
              />
            ))}
          </LineChart>
        ) : (
          <AreaChart data={chartData}>
            <ChartTooltip
              content={
                <ChartTooltipContent
                  formatter={(value, name, item) => (
                    <div className="flex flex-1 justify-between items-center gap-4">
                      <div className="flex items-center gap-2">
                        <div
                          className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                          style={{ backgroundColor: item.color }}
                        />
                        <span className="text-muted-foreground truncate max-w-[150px]">
                          {name}
                        </span>
                      </div>
                      <span className="font-mono font-medium">
                        {formatter(value as number)}
                      </span>
                    </div>
                  )}
                />
              }
            />
            <XAxis
              dataKey="ts"
              tickFormatter={(value) => dayjs(value).format("HH:mm")}
            />
            <YAxis tickFormatter={formatter} width={70} />
            <defs>
              {series.map((key, idx) => (
                <linearGradient
                  key={key}
                  id={`gradient-${index}-${idx}`}
                  x1="0"
                  y1="0"
                  x2="0"
                  y2="1"
                >
                  <stop
                    offset="5%"
                    stopColor={dynamicConfig[key]?.color}
                    stopOpacity={0.5}
                  />
                  <stop
                    offset="95%"
                    stopColor={dynamicConfig[key]?.color}
                    stopOpacity={0.1}
                  />
                </linearGradient>
              ))}
            </defs>
            {series.map((key, idx) => (
              <Area
                key={key}
                dataKey={key}
                type="basis"
                fill={`url(#gradient-${index}-${idx})`}
                fillOpacity={0.2}
                stroke={dynamicConfig[key]?.color}
              />
            ))}
          </AreaChart>
        )}
      </ChartContainer>
      <div className="text-xs text-gray-600 flex items-center gap-1">
        <InfoIcon size={12} /> {gauge.description}
      </div>
    </MetricCard>
  );
}
