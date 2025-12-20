import type { Metadata } from "next";

import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";
import TiKV from "@/assets/img/tikv.webp";
import Sidebar from "@/components/sidebar";
import Provider from "@/components/provider";
import Cluster from "@/components/cluster";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "TiKV UI by Stream",
  description: "TiKV UI by Stream",
  icons: {
    icon: TiKV.src,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Provider>
          <div className="flex flex-1 flex-row h-screen">
            <Sidebar />
            {children}
          </div>
          <Cluster />
          <Toaster position="bottom-right" />
        </Provider>
      </body>
    </html>
  );
}
