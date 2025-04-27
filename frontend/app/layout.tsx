import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ReactQueryProviders } from "@/providers/react-query.provider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata = {
  title: 'Kith Monitor',
  description: 'Moniteur de produits Kith',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ReactQueryProviders>
          {children}
        </ReactQueryProviders>
      </body>
    </html>
  );
}