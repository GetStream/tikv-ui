import { ScrollArea } from "@/components/ui/scroll-area";
import { ScanItem } from "@/hooks/use-keys";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/esm/styles/prism";

interface KeyDetailsContentProps {
  selectedItem: ScanItem;
  view: "raw" | "parsed";
}

export function KeyDetailsContent({
  selectedItem,
  view,
}: KeyDetailsContentProps) {
  return (
    <ScrollArea className="h-full">
      <div className="text-muted-foreground break-all whitespace-pre-wrap">
        <div className="relative pb-15">
          {view == "parsed" ? (
            <SyntaxHighlighter
              language="json"
              style={vscDarkPlus}
              customStyle={{
                borderRadius: "0.5rem",
                fontSize: "0.875rem",
                background: "transparent",
              }}
            >
              {JSON.stringify(selectedItem.value, null, 2)}
            </SyntaxHighlighter>
          ) : (
            selectedItem.raw_value
          )}
        </div>
      </div>
    </ScrollArea>
  );
}
