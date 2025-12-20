"use client";

import { CheckCircle2Icon } from "lucide-react";
import { useState } from "react";
import { useCluster } from "@/hooks/use-cluster";
import { useKeys } from "@/hooks/use-keys";
import { SwitchClusterDialog } from "@/components/dialogs/switch-cluster";

export default function Cluster() {
  const { loadKeys } = useKeys();
  const { clusters, switchCluster } = useCluster();
  const [open, setOpen] = useState(false);

  const activeCluster = clusters.find((c) => c.active);

  const handleSwitch = async (name: string) => {
    await switchCluster(name);
    loadKeys("", "", true);
    setOpen(false);
  };

  if (!activeCluster && clusters.length === 0) return null;

  return (
    <>
      <div className="fixed bottom-5 right-5" onClick={() => setOpen(true)}>
        <div className="flex items-center gap-2 text-xs opacity-60 hover:opacity-100 transition cursor-pointer bg-background/80 backdrop-blur-sm p-2 rounded-md border border-border shadow-sm">
          <CheckCircle2Icon
            size={15}
            className="text-emerald-500"
            strokeWidth={3}
          />
          <div className="font-mono">
            {activeCluster
              ? activeCluster.pd_addrs.join(",")
              : "No active cluster"}
          </div>
        </div>
      </div>

      <SwitchClusterDialog
        open={open}
        onOpenChange={setOpen}
        clusters={clusters}
        onSwitch={handleSwitch}
      />
    </>
  );
}
