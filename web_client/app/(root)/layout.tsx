import type { Metadata } from "next";
import "../globals.css";
import LeftSideBar from "@/components/shared/LeftSideBar";
import TopBar from "@/components/shared/TopBar";
import { Theme } from "@/context/Theme";
import { Poppins } from 'next/font/google';
import { Toaster } from "@/components/ui/toaster";
import { UserProvider } from "@/context/UserContext";
import ReactQueryProvider from "@/context/ReactQueryProvider";

const poppins = Poppins({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-poppins',
  weight: ['100', '200', '300', '400', '500', '600', '700', '800', '900']
});

export const metadata: Metadata = {
  title: "Aethly",
  description: "Cloud Storage",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <ReactQueryProvider>
      <body className={`${poppins.variable} flex flex-row h-screen overflow-y-hidden m-4`}>
        <UserProvider>
        <Theme>
          <div className="p-3">
            <LeftSideBar />
          </div>

          <div className="flex min-h-screen flex-1 flex-col w-full p-3">
            <div className=""><TopBar /></div>
            <div className="w-full p-4 dark:bg-black bg-black/10 rounded-xl my-3 h-screen overflow-y-auto">{children}</div>
          </div>
        </Theme>
        <Toaster />
        </UserProvider>
      </body>
      </ReactQueryProvider>
    </html>
  );
}
