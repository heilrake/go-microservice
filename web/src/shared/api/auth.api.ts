import type { HTTPDriverLoginRequestPayload, HTTPUserLoginRequestPayload } from '@/shared/libs/contracts'

import { internalClient } from './client'

export const authApi = {
  me: () =>
    internalClient.get<{ userID: string } | null>('/api/auth/me'),

  loginRider: (payload: HTTPUserLoginRequestPayload) =>
    internalClient.post<void>('/api/auth/login', payload),

  loginDriver: (payload: HTTPDriverLoginRequestPayload) =>
    internalClient.post<void>('/api/auth/driver-login', payload),

  logout: () =>
    internalClient.post<void>('/api/auth/logout', {}),
}
