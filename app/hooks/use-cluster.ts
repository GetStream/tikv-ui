"use client";

import { useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
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
  const queryClient = useQueryClient();

  const {
    data: clusters = [],
    isPending: isQueryLoading,
    error: queryError,
    refetch,
  } = useQuery({
    queryKey: ["clusters"],
    queryFn: async () => {
      const response = await fetch(`${API_BASE_URL}/api/clusters`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: ClustersResponse = await response.json();
      return data.clusters;
    },
  });

  const connectMutation = useMutation({
    mutationFn: async ({
      name,
      pdAddrs,
    }: {
      name: string;
      pdAddrs: string[];
    }) => {
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
      return data;
    },
    onSuccess: (data) => {
      toast.success("Connected to cluster successfully", {
        description: `Cluster "${data.name}" is now active`,
      });
      queryClient.invalidateQueries({ queryKey: ["clusters"] });
      queryClient.invalidateQueries({ queryKey: ["metrics"] });
    },
    onError: (err) => {
      const msg =
        err instanceof Error ? err.message : "Failed to connect to cluster";
      toast.error("Connection failed", { description: msg });
    },
  });

  const switchMutation = useMutation({
    mutationFn: async (name: string) => {
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
    },
    onSuccess: (_, name) => {
      toast.success("Switched cluster", {
        description: `Active cluster is now "${name}"`,
      });
      queryClient.invalidateQueries({ queryKey: ["clusters"] });
      queryClient.refetchQueries({ queryKey: ["metrics"] });
    },
    onError: (err) => {
      const msg =
        err instanceof Error ? err.message : "Failed to switch cluster";
      toast.error("Switch failed", { description: msg });
    },
  });

  const loading =
    isQueryLoading || connectMutation.isPending || switchMutation.isPending;

  const error =
    (queryError as Error)?.message ||
    (connectMutation.error as Error)?.message ||
    (switchMutation.error as Error)?.message ||
    null;

  const { mutateAsync: connectMutate } = connectMutation;
  const { mutateAsync: switchMutate } = switchMutation;

  const connectCluster = useCallback(
    (name: string, pdAddrs: string[]) => connectMutate({ name, pdAddrs }),
    [connectMutate]
  );

  const listClusters = useCallback(async () => {
    await refetch();
  }, [refetch]);

  const switchCluster = useCallback(
    (name: string) => switchMutate(name),
    [switchMutate]
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
