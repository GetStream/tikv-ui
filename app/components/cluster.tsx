import { CheckCircle2Icon, ServerIcon } from "lucide-react";
import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { ClusterInfo } from "@/hooks/use-cluster";

export default function Cluster({
  onChange,
  switchCluster,
  clusters,
}: {
  onChange: () => void;
  switchCluster: (name: string) => Promise<void>;
  clusters: ClusterInfo[];
}) {
  const [open, setOpen] = useState(false);

  const activeCluster = clusters.find((c) => c.active);

  const handleSwitch = async (name: string) => {
    await switchCluster(name);
    onChange?.();
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

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Switch Cluster</DialogTitle>
          </DialogHeader>
          <div className="grid gap-2 py-4">
            {clusters.map((cluster) => (
              <div
                key={cluster.name}
                className={cn(
                  "flex items-center justify-between p-3 rounded-md border cursor-pointer transition-colors",
                  cluster.active
                    ? "bg-accent border-accent-foreground/20"
                    : "hover:bg-muted"
                )}
                onClick={() => handleSwitch(cluster.name)}
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
    </>
  );
}
