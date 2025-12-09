import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { BanIcon, KeyIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import { ScanItem } from "@/hooks/use-keys";

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
}: KeyListProps) {
  return (
    <div className="flex-1 min-w-0 overflow-hidden">
      <div className="p-4">
        <Input
          placeholder="Search keys..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
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
            <div className="space-y-1 w-full">
              {keys.map((item, index) => {
                return (
                  <div
                    key={item.key}
                    onClick={() => setSelectedItem(item)}
                    className={cn(
                      "p-2 rounded hover:bg-primary/30 cursor-pointer text-sm font-mono w-[280px] rounded-md",
                      selectedItem?.key === item.key ? "bg-primary" : ""
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
  );
}
