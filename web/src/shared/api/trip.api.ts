import type { HTTPTripStartResponse } from '@/features/trip';

import type {
  HTTPTripPreviewRequestPayload,
  HTTPTripPreviewResponse,
  HTTPTripStartRequestPayload,
} from '@/shared/libs/contracts';

import { gatewayClient } from './client';

export const tripApi = {
  preview: (payload: HTTPTripPreviewRequestPayload) =>
    gatewayClient.post<HTTPTripPreviewResponse>('/trip/preview', payload),

  start: (payload: HTTPTripStartRequestPayload) =>
    gatewayClient.post<HTTPTripStartResponse>('/trip/start', payload),

  cancel: (payload: { userID: string }) => gatewayClient.post<void>('/trip/cancel', payload),
};
