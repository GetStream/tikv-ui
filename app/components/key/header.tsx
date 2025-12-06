import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import { ScanItem } from "@/hooks/use-keys";
import {
  BookCheckIcon,
  BookDashedIcon,
  KeyIcon,
  TrashIcon,
} from "lucide-react";
import { formatBytes } from "@/lib/utils";

interface KeyDetailsHeaderProps {
  selectedItem: ScanItem | null;
  view: "raw" | "parsed";
  setView: (view: "raw" | "parsed") => void;
  onDelete: () => void;
}

export function KeyDetailsHeader({
  selectedItem,
  view,
  setView,
  onDelete,
}: KeyDetailsHeaderProps) {
  if (!selectedItem) return null;

  return (
    <div className="flex items-center justify-between mb-8 border-b pb-8">
      <div className="font-semibold flex items-center gap-2">
        <KeyIcon size={14} strokeWidth={3} />
        {selectedItem.key}
        <span className="text-xs text-muted-foreground ml-2 font-normal">
          ({formatBytes(selectedItem.raw_value.length)})
        </span>
      </div>

      <div className="flex items-center gap-2">
        <ButtonGroup aria-label="Button group">
          <Button
            variant={view === "parsed" ? "secondary" : "outline"}
            onClick={() => setView("parsed")}
          >
            <BookCheckIcon /> Parsed
          </Button>
          <Button
            variant={view === "raw" ? "secondary" : "outline"}
            onClick={() => setView("raw")}
          >
            <BookDashedIcon /> Raw
          </Button>
        </ButtonGroup>
        <Button variant="destructive" onClick={onDelete}>
          <TrashIcon />
        </Button>
      </div>
    </div>
  );
}
