import { PlatformType, Tier } from '@odigos/ui-kit/types';
import type { OperationContext } from '@odigos/ui-kit/contexts';

export const IS_DEV = process.env.NODE_ENV === 'development';

const isLoopbackHost = typeof window !== 'undefined' ? /^(localhost|127\.0\.0\.1|\[::1\])$/.test(window.location.hostname) : false;

export const IS_LOCAL = IS_DEV && isLoopbackHost;

/** Public go offsets manifest (same URL as `odigos pro update-offsets`). */
export const GO_OFFSETS_PUBLIC_URL = 'https://storage.googleapis.com/odigos-cloud/offset_results_min.json';

/**
 * Initial operation context used while we bootstrap. The real context
 * (platformType/tier/version) is derived by the layout via `useConfig`
 * and propagated through React rerenders.
 */
export const INITIAL_CONTEXT: OperationContext = {
  platformType: PlatformType.K8s,
  tier: Tier.Community,
  version: 'v0.0.0',
};
