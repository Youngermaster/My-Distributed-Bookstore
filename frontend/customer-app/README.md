# Customer App (Frontend)

## Overview

React-based customer-facing web application for the distributed bookstore system.

## Technology Stack

- **Language**: TypeScript
- **Framework**: React 18+
- **UI Library**: shadcn/ui
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Data Fetching**: TanStack Query v5
- **Routing**: React Router v6
- **Build Tool**: Vite
- **Package Manager**: npm

## Project Structure

```
customer-app/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn/ui components
│   │   ├── layout/          # Header, Footer, Sidebar
│   │   ├── books/           # Book-related components
│   │   ├── cart/            # Cart components
│   │   ├── orders/          # Order components
│   │   └── auth/            # Authentication components
│   ├── pages/
│   │   ├── HomePage.tsx
│   │   ├── BookDetailsPage.tsx
│   │   ├── SearchPage.tsx
│   │   ├── CartPage.tsx
│   │   ├── CheckoutPage.tsx
│   │   ├── OrderHistoryPage.tsx
│   │   └── ProfilePage.tsx
│   ├── api/                 # API client (Axios)
│   ├── hooks/               # Custom React hooks
│   ├── store/               # Zustand stores
│   ├── types/               # TypeScript type definitions
│   ├── utils/               # Utility functions
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css
├── public/
├── Dockerfile
├── nginx.conf
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── vite.config.ts
└── README.md
```

## Features

- Book browsing and search
- Shopping cart management
- User authentication (JWT)
- Order placement and tracking
- User profile management
- Responsive design
- Dark mode support
- Book reviews and ratings
- Personalized recommendations

## Getting Started

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm test
```

## Environment Variables

Create a `.env` file:

```bash
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Bookstore
```

## Next Steps

- [ ] Initialize Vite + React + TypeScript project
- [ ] Install and configure shadcn/ui
- [ ] Setup Tailwind CSS
- [ ] Create API client with Axios
- [ ] Setup TanStack Query
- [ ] Setup Zustand stores
- [ ] Create UI components
- [ ] Implement pages
- [ ] Add authentication flow
- [ ] Write tests
- [ ] Configure production build
