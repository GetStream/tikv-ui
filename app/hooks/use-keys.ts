"use client";

import { useCallback, useState, useMemo } from "react";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
import { API_BASE_URL } from "@/lib/utils";

export interface ScanItem {
  key: string;
  value: any;
  raw_value: string;
}

interface ScanResponse {
  items: ScanItem[];
}

export function useKeys() {
  const queryClient = useQueryClient();

  // We track the initial start/end keys for the current "list session" (e.g. search or initial load)
  const [filter, setFilter] = useState({ startKey: "", endKey: "" });

  const getNextKey = (key: string) => {
    if (!key) return "";
    const lastChar = key.slice(-1);
    const prefix = key.slice(0, -1);
    return prefix + String.fromCharCode(lastChar.charCodeAt(0) + 1);
  };

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetching,
    // error as unknown as Error to match common simplified signature, or handle properly
    error,
    refetch,
  } = useInfiniteQuery({
    queryKey: ["keys", filter],
    queryFn: async ({ pageParam }) => {
      // If pageParam is set, it overrides the filter's startKey for pagination
      const currentStartKey = pageParam ?? filter.startKey;

      const response = await fetch(`${API_BASE_URL}/api/raw/scan`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          start_key: currentStartKey,
          end_key: filter.endKey,
          limit: 100,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: ScanResponse = await response.json();
      return data.items || [];
    },
    initialPageParam: null as string | null,
    getNextPageParam: (lastPage, allPages) => {
      if (!lastPage || lastPage.length < 100) {
        return undefined;
      }
      const lastItem = lastPage[lastPage.length - 1];
      // We use the raw key here. The existing logic slices the first item if it matches startKey
      // or we can just pass the last key and let the de-duplication happen in display or here.
      // However, the previous implementation did:
      // loadKeys(last.key, ...)
      // and then removed the first item if it matched.
      // To keep it simple, we pass last.key.
      return lastItem.key;
    },
    // Don't refetch on window focus as it resets the list state weirdly for infinite lists sometimes
    refetchOnWindowFocus: false,
  });

  // Flatten the pages into a single array
  const keys = useMemo(() => {
    if (!data?.pages) return [];
    return data.pages.flatMap((page, pageIndex) => {
      // For pages after the first one, we might want to slice the first item if it overlaps
      // But based on `getNextPageParam` returning `lastItem.key`, the next page will start with that key.
      // The default behavior of TiKV scan (usually inclusive start) implies overlap.
      // Previous implementation:
      // const itemsToAdd = newItems.length > 0 && newItems[0].key === startKey ? newItems.slice(1) : newItems;
      if (pageIndex > 0 && page.length > 0) {
        // We can't easily check against "previous page last key" here without looking it up
        // But we know 'pageParam' for this page was the last key of previous page.
        // Let's heuristically checks if this page's first key matches the last key of the previous page.
        const prevPage = data.pages[pageIndex - 1];
        const lastKeyOfPrev = prevPage[prevPage.length - 1]?.key;
        if (page[0].key === lastKeyOfPrev) {
          return page.slice(1);
        }
      }
      return page;
    });
  }, [data]);

  const addMutation = useMutation({
    mutationFn: async ({ key, value }: { key: string; value: string }) => {
      const response = await fetch(`${API_BASE_URL}/api/raw/put`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ key, value }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    },
    onSuccess: (_, { key }) => {
      toast.success("Key added successfully", {
        description: `"${key}" has been added`,
      });
      // Invalidate query to refresh list
      queryClient.invalidateQueries({ queryKey: ["keys"] });
    },
    onError: (err) => {
      toast.error("Failed to add key", {
        description: err instanceof Error ? err.message : "Unknown error",
      });
      throw err; // Propagate for UI handling if needed
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (key: string) => {
      const response = await fetch(`${API_BASE_URL}/api/raw/delete`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ key }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    },
    onSuccess: (_, key) => {
      toast.success("Key deleted successfully", {
        description: `"${key}" has been removed`,
      });
      queryClient.invalidateQueries({ queryKey: ["keys"] });
    },
    onError: (err) => {
      toast.error("Failed to delete key", {
        description: err instanceof Error ? err.message : "Unknown error",
      });
    },
  });

  // Preserve the API signature expected by consumers
  const loadKeys = useCallback(
    async (
      startKey: string,
      endKey: string = "",
      replaceList: boolean = false
    ) => {
      if (replaceList) {
        // New search or initial load
        setFilter({ startKey, endKey });
        refetch();
        // The effect of setFilter will convert to a new query key and trigger fetch
      } else {
        // Load more
        await fetchNextPage();
      }
    },
    [fetchNextPage]
  );

  return {
    keys,
    loading: isFetching,
    error: error ? (error as Error).message : null,
    hasMore: !!hasNextPage,
    loadKeys,
    addKey: async (key: string, value: string) =>
      addMutation.mutateAsync({ key, value }),
    deleteKey: async (key: string) => deleteMutation.mutateAsync(key),
    getNextKey,
  };
}
