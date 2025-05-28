import type { Metadata } from "next";
import "./globals.css";
import { Toaster } from "sonner";
import { Inter, Poppins } from "next/font/google";

// const inter = Inter({
//   subsets: ["latin"],
//   variable: "--font-inter",
//   display: "swap",
// });

const poppins = Poppins({
    subsets: ["latin"],
    variable: "--font-poppins",
    display: "swap",
    weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
});
export const metadata: Metadata = {
  title: "ClipIt",
  description: "Just the moments that matter",
  openGraph: {
    title: "ClipIt",
    description: "Just the moments that matter",
    type: "website",
    siteName: "ClipIt",
    url: "https://clipit.avadhootsmart.xyz",
  },
  twitter: {
    title: "ClipIt",
    description: "Just the moments that matter",
    card: "summary_large_image",
    creator: "@avadhoot_smart",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${poppins.variable} antialiased`}>
        <main>{children}</main>
        <Toaster richColors position="top-center" />
      </body>
    </html>
  );
}
