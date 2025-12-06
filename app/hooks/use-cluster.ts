"use client";

import { useCallback, useState } from "react";
import { toast } from "sonner";
import { API_BASE_URL } from "@/lib/utils";

export interface ClusterInfo {
  name: string;
  cluster_id: string;
  pd_addrs: string[];
  active: boolean;
}

export interface ClustersResponse {
  clusters: ClusterInfo[];
}

export function useCluster() {
  const [clusters, setClusters] = useState<ClusterInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const connectCluster = useCallback(
    async (name: string, pdAddrs: string[]) => {
      try {
        setLoading(true);
        setError(null);
        const response = await fetch(`${API_BASE_URL}/api/clusters/connect`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name,
            pd_addrs: pdAddrs,
          }),
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(
            errorData.error || `HTTP error! status: ${response.status}`
          );
        }

        const data: ClusterInfo = await response.json();

        toast.success("Connected to cluster successfully", {
          description: `Cluster "${data.name}" is now active`,
        });

        return data;
      } catch (err) {
        const msg =
          err instanceof Error ? err.message : "Failed to connect to cluster";
        setError(msg);
        toast.error("Connection failed", { description: msg });
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const listClusters = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/clusters`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: ClustersResponse = await response.json();
      setClusters(data.clusters);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to list clusters");
    } finally {
      setLoading(false);
    }
  }, []);

  const switchCluster = useCallback(
    async (name: string) => {
      try {
        setLoading(true);
        const response = await fetch(`${API_BASE_URL}/api/clusters/switch`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name,
          }),
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(
            errorData.error || `HTTP error! status: ${response.status}`
          );
        }

        toast.success("Switched cluster", {
          description: `Active cluster is now "${name}"`,
        });

        // Refresh list to update active status
        await listClusters();
      } catch (err) {
        const msg =
          err instanceof Error ? err.message : "Failed to switch cluster";
        setError(msg);
        toast.error("Switch failed", { description: msg });
      } finally {
        setLoading(false);
      }
    },
    [listClusters]
  );

  return {
    clusters,
    loading,
    error,
    connectCluster,
    listClusters,
    switchCluster,
  };
}
