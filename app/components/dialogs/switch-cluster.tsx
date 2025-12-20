"use client";

import { CheckCircle2Icon, ServerIcon } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { ClusterInfo } from "@/hooks/use-cluster";

interface SwitchClusterDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  clusters: ClusterInfo[];
  onSwitch: (name: string) => void;
}

export function SwitchClusterDialog({
  open,
  onOpenChange,
  clusters,
  onSwitch,
}: SwitchClusterDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Switch Cluster</DialogTitle>
        </DialogHeader>
        <div className="grid gap-2 py-4">
          {clusters.map((cluster) => (
            <div
              key={cluster.name}
              className={cn(
                "flex items-center justify-between p-3 rounded-md border border-secondary cursor-pointer transition-colors",
                cluster.active
                  ? "bg-black/20 border-transparent"
                  : "hover:bg-black/10 hover:border-transparent"
              )}
              onClick={() => onSwitch(cluster.name)}
            >
              <div className="flex items-center gap-3">
                <ServerIcon size={18} className="opacity-70" />
                <div className="flex flex-col">
                  <span className="text-sm font-medium">{cluster.name}</span>
                  <span className="text-xs text-muted-foreground">
                    {cluster.pd_addrs.join(", ")}
                  </span>
                </div>
              </div>
              {cluster.active && (
                <CheckCircle2Icon size={16} className="text-emerald-500" />
              )}
            </div>
          ))}
          {clusters.length === 0 && (
            <div className="text-center text-sm text-muted-foreground py-4">
              No clusters found.
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

