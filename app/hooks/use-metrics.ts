"use client";

import { API_BASE_URL } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";

type MetricsResponse = {
  pd: {
    count: number;
    stores: {
      store: {
        id: number;
        address: string;
        state: number;
        last_heartbeat: number;
        state_name: string;
        status_address: string;
      };
      status: {
        capacity: string;
        available: string;
        used_size: string;
        leader_count: number;
        region_count: number;
        start_ts: string;
        last_heartbeat_ts: string;
        uptime: string;
      };
    }[];
  };
  tikv: {
    gauges: {
      name: string;
      labels: Record<string, string>;
      unit: string;
      points: {
        ts: number;
        value: number;
      }[];
    }[];
  };
};
export function useMetrics() {
  return useQuery({
    queryKey: ["metrics"],
    queryFn: async () => {
      const response = await fetch(`${API_BASE_URL}/api/metrics`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: MetricsResponse = await response.json();
      return data;
    },
    refetchInterval: 1000 * 10,
  });
}
