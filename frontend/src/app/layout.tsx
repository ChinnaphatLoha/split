import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Split — Smart Group Expense Manager',
  description: 'Track shared expenses, split bills effortlessly, and settle debts with the optimal payment plan. Built for friends, roommates, and teams.',
  keywords: ['expense tracker', 'split bills', 'group expenses', 'debt settlement', 'money management'],
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800;900&display=swap"
          rel="stylesheet"
        />
      </head>
      <body style={{ margin: 0, backgroundColor: '#0A0E1A' }}>
        {children}
      </body>
    </html>
  );
}
