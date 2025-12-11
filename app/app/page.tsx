"use client";

import { MousePointerClickIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { DeleteKeyDialog } from "@/components/dialogs/deletekey";

import { useKeys, ScanItem } from "@/hooks/use-keys";
import Sidebar from "@/components/sidebar";
import { KeyList } from "@/components/key/list";
import { KeyDetailsContent } from "@/components/key/content";
import { KeyDetailsHeader } from "@/components/key/header";
import Cluster from "@/components/cluster";
import { useCluster } from "@/hooks/use-cluster";

export default function Home() {
  const {
    keys,
    loading,
    error,
    hasMore,
    loadKeys,
    addKey,
    deleteKey,
    getNextKey,
  } = useKeys();

  const [selectedItem, setSelectedItem] = useState<ScanItem | null>(null);
  const [view, setView] = useState<"raw" | "parsed">("parsed");
  const [searchQuery, setSearchQuery] = useState("");
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const { clusters, listClusters, switchCluster } = useCluster();

  const [pdAddrs, setPdAddrs] = useState("");

  useEffect(() => {
    const savedAddrs = localStorage.getItem("tikv_pd_addrs");
    if (savedAddrs) {
      setPdAddrs(savedAddrs);
    }
  }, []);

  // Initial load
  useEffect(() => {
    listClusters();
    loadKeys("", "", true);
  }, [loadKeys, listClusters]);

  // Debounced search effect
  useEffect(() => {
    const timer = setTimeout(() => {
      if (searchQuery.trim()) {
        const endKey = getNextKey(searchQuery);
        loadKeys(searchQuery, endKey, true);
      } else {
        loadKeys("", "", true);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery, loadKeys]);

  const handleDeleteWrapper = async () => {
    if (!selectedItem) return;
    await deleteKey(selectedItem.key);
    setSelectedItem(null);

    if (searchQuery) {
      const endKey = getNextKey(searchQuery);
      loadKeys(searchQuery, endKey, true);
    } else {
      loadKeys("", "", true);
    }
  };

  return (
    <div className="flex flex-1 flex-row h-screen">
      <div className="flex flex-1 flex-row border-r border-border max-w-sm">
        <Sidebar
          addKey={addKey}
          loadKeys={loadKeys}
          getNextKey={getNextKey}
          searchQuery={searchQuery}
          listClusters={listClusters}
        />
        <KeyList
          keys={keys}
          loading={loading}
          error={error}
          hasMore={hasMore}
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          selectedItem={selectedItem}
          setSelectedItem={setSelectedItem}
          loadKeys={loadKeys}
          getNextKey={getNextKey}
        />
      </div>
      <div className="flex-1 p-8 overflow-hidden">
        <KeyDetailsHeader
          selectedItem={selectedItem}
          view={view}
          setView={setView}
          onDelete={() => setDeleteDialogOpen(true)}
        />

        {selectedItem ? (
          <KeyDetailsContent selectedItem={selectedItem} view={view} />
        ) : (
          <div className="flex flex-col items-center justify-center gap-2 h-screen opacity-50">
            <MousePointerClickIcon size={130} className="opacity-50" />
            <div className="text-sm">Select a key to view details</div>
          </div>
        )}
        <Cluster
          onChange={() => loadKeys("", "", true)}
          switchCluster={switchCluster}
          clusters={clusters}
        />
      </div>

      <DeleteKeyDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        itemKey={selectedItem?.key}
        onConfirm={handleDeleteWrapper}
      />
    </div>
  );
}
