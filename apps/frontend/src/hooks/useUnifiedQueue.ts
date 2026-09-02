import { useMemo } from 'react';
import type { CaseItemDTO } from '../api/types';
import type { LiveAction } from './useLiveActions';

export interface UnifiedQueueItem {
  case_id: string;
  channel: string;
  timestamp: number;
  amount_at_risk: number;
  customer_id: string;
  is_cold: boolean;
}

export const useUnifiedQueue = (cases: CaseItemDTO[] | null | undefined = [], liveActions: LiveAction[] = []): UnifiedQueueItem[] => {
  return useMemo(() => {
    const displayedCasesMap = new Map<string, UnifiedQueueItem>();

    // 1. Cold Start Cases from API
    if (cases && Array.isArray(cases)) {
      cases.forEach((c) => {
        const channel = c.latest_action_channel?.Valid ? c.latest_action_channel.String : c.status;
        if (channel === 'silent_retry') return; // Ignore silent_retry

        displayedCasesMap.set(c.id, {
          case_id: c.id,
          channel: channel,
          timestamp: new Date(c.created_at).getTime(),
          amount_at_risk: c.amount_at_risk,
          customer_id: '',
          is_cold: true
        });
      });
    }

    // 2. Stream Overrides from SSE
    if (liveActions && Array.isArray(liveActions)) {
      liveActions.forEach((a) => {
        if (a.channel === 'silent_retry') return; // Ignore silent_retry

        const existing = displayedCasesMap.get(a.case_id) || ({} as Partial<UnifiedQueueItem>);
        displayedCasesMap.set(a.case_id, {
          case_id: a.case_id,
          channel: a.channel,
          customer_id: a.customer_id,
          timestamp: a.timestamp || existing.timestamp || Date.now(),
          amount_at_risk: existing.amount_at_risk || 0,
          is_cold: false
        });
      });
    }

    // 3. Sort by most recent
    return Array.from(displayedCasesMap.values()).sort((a, b) => b.timestamp - a.timestamp);
  }, [cases, liveActions]);
};
