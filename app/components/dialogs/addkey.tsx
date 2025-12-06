"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { cn } from "@/lib/utils";

const formSchema = z.object({
  key: z.string().min(1, "Key is required"),
  value: z.string().min(1, "Value is required"),
});

type FormValues = z.infer<typeof formSchema>;

interface AddKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAdd: (key: string, value: string) => Promise<void>;
}

export function AddKeyDialog({ open, onOpenChange, onAdd }: AddKeyDialogProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  });

  const onSubmit = async (data: FormValues) => {
    try {
      await onAdd(data.key, data.value);
      reset();
      onOpenChange(false);
    } catch {
      // Error handled by parent
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Record</DialogTitle>
          <DialogDescription>
            Create a new key-value pair in the database.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <div className="flex flex-col gap-2">
              <label htmlFor="key" className="text-sm font-medium">
                Key
              </label>
              <Input
                id="key"
                placeholder="key"
                className={cn(errors.key && "border-destructive")}
                {...register("key")}
              />
              {errors.key && (
                <p className="text-xs text-destructive">{errors.key.message}</p>
              )}
            </div>
            <div className="flex flex-col gap-2">
              <label htmlFor="value" className="text-sm font-medium">
                Value
              </label>
              <textarea
                id="value"
                placeholder="value"
                className={cn(
                  "flex min-h-[150px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
                  errors.value && "border-destructive"
                )}
                {...register("value")}
              />
              {errors.value && (
                <p className="text-xs text-destructive">
                  {errors.value.message}
                </p>
              )}
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              size="lg"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" size="lg" disabled={isSubmitting}>
              {isSubmitting ? "Adding..." : "Add"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
