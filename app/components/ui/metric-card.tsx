import { cn } from "@/lib/utils";

interface MetricCardProps {
  children: React.ReactNode;
  className?: string;
}

export function MetricCard({ children, className }: MetricCardProps) {
  return (
    <div
      className={cn(
        "rounded-lg p-4 bg-menu hover:bg-secondary/20 transition border border-transparent hover:border-secondary",
        className
      )}
    >
      {children}
    </div>
  );
}

