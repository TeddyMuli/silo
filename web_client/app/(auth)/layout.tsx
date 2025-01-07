import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "../globals.css";
import { Toaster } from "@/components/ui/toaster";
import { Theme } from "@/context/Theme";
import Logo from "@/components/shared/Logo";
import ReactQueryProvider from "@/context/ReactQueryProvider";

const poppins = Poppins({
  subsets: ['latin'],
  display: 'swap',
  weight: ['100', '200', '300', '400', '500', '600', '700', '800', '900']
});

export const metadata: Metadata = {
  title: "Silo",
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
        <body className={`${poppins.variable} m-4 bg-black`}>
          <Theme>
            <Logo />
            {children}
            <Toaster />
          </Theme>
        </body>
      </ReactQueryProvider>
    </html>
  );
}
