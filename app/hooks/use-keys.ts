"use client";

import { useCallback, useState } from "react";
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
  const [keys, setKeys] = useState<ScanItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);

  const getNextKey = (key: string) => {
    if (!key) return "";
    const lastChar = key.slice(-1);
    const prefix = key.slice(0, -1);
    return prefix + String.fromCharCode(lastChar.charCodeAt(0) + 1);
  };

  const loadKeys = useCallback(
    async (
      startKey: string,
      endKey: string = "",
      replaceList: boolean = false
    ) => {
      try {
        setLoading(true);
        const response = await fetch(`${API_BASE_URL}/api/raw/scan`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            start_key: startKey,
            end_key: endKey,
            limit: 100,
          }),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data: ScanResponse = await response.json();
        const newItems = data.items || [];

        if (newItems.length < 100) {
          setHasMore(false);
        } else {
          setHasMore(true);
        }

        if (replaceList) {
          setKeys(newItems);
        } else {
          const itemsToAdd =
            newItems.length > 0 && newItems[0].key === startKey
              ? newItems.slice(1)
              : newItems;
          setKeys((prev) => [...prev, ...itemsToAdd]);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch keys");
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const addKey = async (key: string, value: string) => {
    const response = await fetch(`${API_BASE_URL}/api/raw/put`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        key,
        value,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    toast.success("Key added successfully", {
      description: `"${key}" has been added`,
    });
  };

  const deleteKey = async (key: string) => {
    const response = await fetch(`${API_BASE_URL}/api/raw/delete`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        key,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    setKeys((prev) => prev.filter((k) => k.key !== key));

    toast.success("Key deleted successfully", {
      description: `"${key}" has been removed`,
    });
  };

  return {
    keys,
    loading,
    error,
    hasMore,
    loadKeys,
    addKey,
    deleteKey,
    getNextKey,
  };
}
