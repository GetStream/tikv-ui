"use client";

import { useState } from "react";
import { Button } from "./ui/button";
import { ListIcon, ServerIcon, SettingsIcon } from "lucide-react";
import Image from "next/image";
import TiKV from "@/assets/img/tikv.webp";
import { SettingsDialog } from "./dialogs/settings";
import { usePathname } from "next/navigation";
import Link from "next/link";
import { useCluster } from "@/hooks/use-cluster";
import { useKeys } from "@/hooks/use-keys";

export default function Sidebar() {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const path = usePathname();
  const { listClusters } = useCluster();
  const { loadKeys } = useKeys();
  return (
    <>
      <div className="p-4 flex flex-col justify-between items-center bg-menu">
        <a href="">
          <Image src={TiKV.src} alt="Logo" width={32} height={32} />
        </a>
        <div className="flex flex-col items-center gap-2">
          <Button
            asChild
            variant={path === "/" ? "secondary" : "ghost"}
            size="icon-lg"
            className="rounded-full"
          >
            <Link href="/">
              <ListIcon
                stroke={path === "/" ? "var(--color-primary)" : "currentColor"}
              />
            </Link>
          </Button>
          <Button
            asChild
            variant={path === "/metrics" ? "secondary" : "ghost"}
            size="icon-lg"
            className="rounded-full"
          >
            <Link href="/metrics">
              <ServerIcon
                stroke={
                  path === "/metrics" ? "var(--color-primary)" : "currentColor"
                }
              />
            </Link>
          </Button>
          <Button
            variant="ghost"
            size="icon-lg"
            className="rounded-full"
            onClick={() => setSettingsOpen(true)}
          >
            <SettingsIcon />
          </Button>
        </div>
      </div>
      <SettingsDialog
        open={settingsOpen}
        onOpenChange={setSettingsOpen}
        onChange={() => {
          listClusters();
          loadKeys("", "", true);
        }}
      />
    </>
  );
}
