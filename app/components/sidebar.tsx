import { useState } from "react";
import { Button } from "./ui/button";
import { PlusIcon, SettingsIcon } from "lucide-react";
import Image from "next/image";
import TiKV from "@/assets/img/tikv.webp";
import { SettingsDialog } from "./dialogs/settings";
import { AddKeyDialog } from "./dialogs/addkey";

export default function Sidebar({
  addKey,
  loadKeys,
  getNextKey,
  searchQuery,
}: any) {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [addDialogOpen, setAddDialogOpen] = useState(false);

  const handleAddWrapper = async (key: string, value: string) => {
    await addKey(key, value);
    // Refresh list
    if (searchQuery) {
      const endKey = getNextKey(searchQuery);
      loadKeys(searchQuery, endKey, true);
    } else {
      loadKeys("", "", true);
    }
  };
  return (
    <>
      <div className="p-4 border-r border-border flex flex-col justify-between items-center">
        <a href="">
          <Image src={TiKV.src} alt="Logo" width={32} height={32} />
        </a>
        <div className="flex flex-col items-center gap-2">
          <Button
            variant="ghost"
            size="icon-lg"
            className="rounded-full"
            onClick={() => setAddDialogOpen(true)}
          >
            <PlusIcon />
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
      <SettingsDialog open={settingsOpen} onOpenChange={setSettingsOpen} />

      <AddKeyDialog
        open={addDialogOpen}
        onOpenChange={setAddDialogOpen}
        onAdd={handleAddWrapper}
      />
    </>
  );
}
