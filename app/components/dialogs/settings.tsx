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
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { cn } from "@/lib/utils";
import { useCluster } from "@/hooks/use-cluster";

const formSchema = z.object({
  pdAddrs: z.string().min(1, "PD Addresses are required"),
});

type FormValues = z.infer<typeof formSchema>;

interface SettingsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function SettingsDialog({ open, onOpenChange }: SettingsDialogProps) {
  const { connectCluster } = useCluster();

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      pdAddrs: "",
    },
  });

  useEffect(() => {
    const savedAddrs = localStorage.getItem("tikv_pd_addrs");
    if (savedAddrs) {
      setValue("pdAddrs", savedAddrs);
    }
  }, [setValue]);

  const onSubmit = async (data: FormValues) => {
    try {
      const addrs = data.pdAddrs
        .split(",")
        .map((a) => a.trim())
        .filter(Boolean);

      const name = "cluster-" + Date.now();

      await connectCluster(name, addrs);

      localStorage.setItem("tikv_pd_addrs", data.pdAddrs);
      onOpenChange(false);
    } catch (error) {
      // Handled by hook toast
      console.error(error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Settings</DialogTitle>
          <DialogDescription>
            Configure your TiKV cluster connection settings.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <div className="flex flex-col gap-2">
              <label htmlFor="pd-addrs" className="text-sm font-medium">
                PD Addresses
              </label>
              <Input
                id="pd-addrs"
                placeholder="127.0.0.1:2379,127.0.0.1:23781..."
                className={cn(errors.pdAddrs && "border-destructive")}
                {...register("pdAddrs")}
              />
              {errors.pdAddrs && (
                <p className="text-xs text-destructive">
                  {errors.pdAddrs.message}
                </p>
              )}
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" disabled={isSubmitting} size="lg">
              {isSubmitting ? "Connecting..." : "Save"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
