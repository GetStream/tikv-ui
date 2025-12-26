"use client";

import { useMetrics } from "@/hooks/use-metrics";
import {
  CircleXIcon,
  CircleCheckIcon,
  ServerIcon,
  ClockIcon,
  HeartIcon,
  HardDriveIcon,
  ChartSplineIcon,
  ServerCrashIcon,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { MetricCard } from "@/components/ui/metric-card";
import { Spinner } from "@/components/ui/spinner";
import { GaugeChart } from "@/components/metrics/gauge-chart";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { parseBytes, formatUptime } from "@/lib/utils";

dayjs.extend(relativeTime);

export default function MetricsPage() {
  const { isLoading, data } = useMetrics();
  return (
    <div className="p-5 w-full">
      <div className="flex items-center gap-2 text-xl font-semibold mb-8 pb-4 border-b">
        <ServerIcon size={18} /> Metrics
      </div>
      {isLoading ? (
        <div className="flex flex-col gap-2 text-gray-500 items-center justify-center h-full mt-[-5%]">
          <div>
            <Spinner className="size-5" />
          </div>
          <div className="text-xs uppercase">loading metrics...</div>
        </div>
      ) : (
        <ScrollArea className="h-full">
          <div className="text-md font-semibold text-gray-500 mb-4">
            PD Stores ({data?.pd?.count})
          </div>
          <div className="grid gap-4 grid-cols-1 md:grid-cols-3 lg:grid-cols-4 3xl:grid-cols-5 mb-20">
            {data?.pd?.stores.map((store) => (
              <MetricCard key={store.store.id}>
                <div className="border-b pb-5 mb-2 relative">
                  <div className="flex items-center gap-2">
                    <div>
                      <div className="text-sm font-semibold">
                        Node: #{store.store.id}
                      </div>
                      <div className="text-xs text-gray-500 font-mono flex items-center gap-1">
                        <ServerCrashIcon size={12} />
                        {store.store.address}
                      </div>
                      <div className="text-xs text-gray-500 font-mono flex items-center gap-1">
                        <ChartSplineIcon size={12} />
                        {store.store.status_address}
                      </div>
                    </div>
                  </div>

                  {store.store.state_name.toLocaleLowerCase() === "up" ? (
                    <Badge
                      variant="outline"
                      className="absolute top-0 right-0 text-green-300/50 border-green-300/50 font-semibold flex items-center gap-1"
                    >
                      <CircleCheckIcon size={13} />
                      Online
                    </Badge>
                  ) : (
                    <Badge
                      variant="outline"
                      className="absolute top-0 right-0 text-red-400/60 border-red-400/60 font-semibold flex items-center gap-1"
                    >
                      <CircleXIcon size={13} />
                      Offline
                    </Badge>
                  )}
                </div>
                <div className="flex items-center justify-between gap-2 border-b p-5 mb-2">
                  <div className="flex items-center justify-center gap-2">
                    <ClockIcon />
                    <div>
                      <div className="text-[0.6rem] text-gray-500 uppercase">
                        uptime
                      </div>
                      <div
                        className="text-gray-200/90"
                        title={store.status.uptime}
                      >
                        {formatUptime(store.status.uptime)}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center justify-center gap-2">
                    <HeartIcon />
                    <div>
                      <div className="text-[0.6rem] text-gray-500 text-center uppercase">
                        heartbeat
                      </div>
                      <div className="text-gray-200/90">
                        {dayjs(store.store.last_heartbeat / 1000000).format(
                          "HH:mm:ss"
                        )}
                      </div>
                    </div>
                  </div>
                </div>
                <div className="mt-5">
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex items-center justify-center gap-1 text-gray-500 text-sm">
                      <HardDriveIcon size={16} /> Storage
                    </div>
                    <div className="text-sm text-gray-500">
                      {store.status.used_size} / {store.status.capacity}
                    </div>
                  </div>
                  <div className="rounded-lg bg-gray-600 relative my-2 h-2 overflow-hidden">
                    <div
                      className="absolute top-0 left-0 w-full h-full bg-green-500/50"
                      style={{
                        width: `${(
                          (parseBytes(store.status.used_size) /
                            parseBytes(store.status.capacity)) *
                          100
                        ).toFixed(2)}%`,
                      }}
                    ></div>
                  </div>
                  <div className="flex items-center justify-between gap-2">
                    <div className="text-sm text-gray-500 ">
                      {store.status.available} available
                    </div>
                    <div className="text-sm text-green-800 font-semibold ">
                      {(
                        (parseBytes(store.status.used_size) /
                          parseBytes(store.status.capacity)) *
                        100
                      ).toFixed(2)}
                      % used
                    </div>
                  </div>
                </div>
              </MetricCard>
            ))}
          </div>
          <div className="text-md font-semibold text-gray-500 mb-4">
            TiKV Overview
          </div>
          <div className="grid gap-4 pb-30 grid-cols-1 md:grid-cols-3 3xl:grid-cols-4 ">
            {data?.tikv?.gauges.map((gauge: any, i: number) => (
              <GaugeChart
                key={gauge.name + i}
                gauge={gauge}
                index={i}
                stores={data?.pd?.stores || []}
              />
            ))}
          </div>
        </ScrollArea>
      )}
    </div>
  );
}
