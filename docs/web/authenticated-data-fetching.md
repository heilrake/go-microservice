# Authenticated Data Fetching Architecture

Next.js App Router + HttpOnly Cookie + SWR Hydration Pattern

---

## Problem

User-authenticated data (e.g. profile, cars, orders, subscriptions, settings, etc.)
must be known **before UI render**.

Fetching such data on the client using:

- `useEffect`
- `useSWR` for initial load
- Zustand auth tokens
- client-side API calls with Authorization headers

causes:

- UI loading flicker
- hydration mismatch
- race conditions
- double fetching during navigation
- layout shift
- inconsistent auth state

---

## Architectural Rule

All authenticated domain data MUST follow:

```
SSR First → SWR Hydration → Client Mutations
```

Initial user data must NEVER be fetched on the client.

---

## Correct Data Flow

### 1. Server Layout Fetch (Initial Truth Source)

Authenticated data must be fetched in:

- Server Components
- Layouts
- Route-level async wrappers

Server fetch must read auth token from:

- HttpOnly cookie

Example:

```tsx id="p8o1i6"
const profile = await getProfileServer();
const orders = await getOrdersServer();
```

This ensures:

- auth-aware rendering
- no loading flicker
- no client race conditions

````

---

### 2. SWR Cache Hydration

Initial server response must be injected into SWR cache using fallback:

```tsx id="jtr43a"
<SWRConfig
  value={{
    fallback: {
      '/api/profile': profile,
      '/api/orders': orders,
    },
    revalidateOnMount: false,
  }}
>
````

This guarantees:

- no fetch after mount
- consistent state across navigation
- instant UI availability

---

### 3. Client Hooks (Read-Only for Initial State)

Client hooks must only read hydrated SWR cache:

```ts id="9ib2r0"
const { data } = useSWR('/api/profile');
```

Client hooks must NOT:

- manage auth token
- fetch initial authenticated data
- expose initial loading state
- rely on Zustand for auth

````

---

### 4. Mutations via Route Handlers

Client must call internal Next.js Route Handlers:

```ts id="nvk7em"
await fetch('/api/profile', { method: 'POST' });
````

After mutation, SWR must revalidate:

```ts id="2k7z5g"
await mutate('/api/profile');
```

DO NOT:

```ts id="p7k1v1"
router.refresh();
```

Mutations must update UI via SWR cache sync.

---

## SWR Usage Policy

SWR is used for:

- background revalidation
- optimistic updates
- post-mutation sync
- shared client cache

SWR is NOT used for:

- initial authenticated data fetch
- auth state detection
- user existence checks

---

## Architecture Summary

```
Server Layout
  → fetch authenticated data
  → hydrate SWR fallback

Client Mount
  → read SWR cache
  → no loading
  → no network

Mutation
  → /api route handler
  → mutate()
  → background revalidate
```

All authenticated resources must follow this pattern.
