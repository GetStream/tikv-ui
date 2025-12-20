import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { BanIcon, KeyIcon, PlusIcon, SearchIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import { ScanItem } from "@/hooks/use-keys";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "../ui/input-group";
import { useState } from "react";
import { AddKeyDialog } from "../dialogs/addkey";

interface KeyListProps {
  keys: ScanItem[];
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  selectedItem: ScanItem | null;
  setSelectedItem: (key: ScanItem | null) => void;
  loadKeys: (
    startKey: string,
    endKey?: string,
    replaceList?: boolean
  ) => Promise<void>;
  getNextKey: (key: string) => string;
  addKey: (key: string, value: string) => Promise<void>;
}

export function KeyList({
  keys,
  loading,
  error,
  hasMore,
  searchQuery,
  setSearchQuery,
  selectedItem,
  setSelectedItem,
  loadKeys,
  getNextKey,
  addKey,
}: KeyListProps) {
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
      <div className="flex-1 min-w-0 overflow-hidden">
        <div className="p-4">
          <InputGroup>
            <InputGroupAddon>
              <SearchIcon />
            </InputGroupAddon>
            <InputGroupInput
              placeholder="Search keys..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <InputGroupAddon align="inline-end">
              <InputGroupButton
                variant="secondary"
                size="icon-sm"
                className="rounded-full"
                onClick={() => setAddDialogOpen(true)}
              >
                <PlusIcon />
              </InputGroupButton>
            </InputGroupAddon>
          </InputGroup>
        </div>
        <ScrollArea className="h-full">
          <div className="p-4">
            {loading && keys.length === 0 && (
              <div className="text-muted-foreground flex items-center gap-2 w-full justify-center">
                <Spinner /> Loading keys...
              </div>
            )}
            {error && <div className="text-destructive">Error: {error}</div>}
            {!loading && !error && keys.length === 0 && (
              <div className="text-muted-foreground flex items-center gap-2 w-full justify-center">
                <BanIcon size={16} />
                No keys found
              </div>
            )}
            {keys.length > 0 && (
              <div className="w-full">
                {keys.map((item, index) => {
                  return (
                    <div
                      key={item.key + index}
                      onClick={() => setSelectedItem(item)}
                      title={item.key}
                      className={cn(
                        "p-2 py-3 rounded hover:bg-secondary cursor-pointer text-sm font-mono w-[355px] rounded-md",
                        selectedItem?.key === item.key
                          ? "bg-primary text-white"
                          : ""
                      )}
                    >
                      <div className="flex items-center gap-2">
                        <KeyIcon size={12} />
                        <div className="truncate w-full">{item.key}</div>
                      </div>
                    </div>
                  );
                })}

                {hasMore ? (
                  <div className="pt-2 pb-20">
                    <Button
                      variant="secondary"
                      className="w-full"
                      size="lg"
                      onClick={() => {
                        const last = keys[keys.length - 1];
                        if (last) {
                          const endKey = searchQuery
                            ? getNextKey(searchQuery)
                            : "";
                          loadKeys(last.key, endKey, false);
                        }
                      }}
                      disabled={loading}
                    >
                      {loading ? (
                        <>
                          <Spinner /> Loading...
                        </>
                      ) : (
                        "Load More"
                      )}
                    </Button>
                  </div>
                ) : (
                  <div className="text-xs text-muted-foreground text-center py-4 pb-20">
                    No more keys
                  </div>
                )}
              </div>
            )}
          </div>
        </ScrollArea>
      </div>
      <AddKeyDialog
        open={addDialogOpen}
        onOpenChange={setAddDialogOpen}
        onAdd={handleAddWrapper}
      />
    </>
  );
}
