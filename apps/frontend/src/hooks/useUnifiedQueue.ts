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

export const useUnifiedQueue = (cases: CaseItemDTO[] = [], liveActions: LiveAction[] = []): UnifiedQueueItem[] => {
  return useMemo(() => {
    const displayedCasesMap = new Map<string, UnifiedQueueItem>();

    // 1. Cold Start Cases from API
    cases.forEach((c) => {
      displayedCasesMap.set(c.id, {
        case_id: c.id,
        // Use the new latest_action_channel if valid, otherwise fallback to the raw status
        channel: c.latest_action_channel?.Valid ? c.latest_action_channel.String : c.status,
        timestamp: new Date(c.created_at).getTime(),
        amount_at_risk: c.amount_at_risk,
        customer_id: '', // Cases API doesn't expose customer_id in list yet
        is_cold: true
      });
    });

    // 2. Stream Overrides from SSE
    liveActions.forEach((a) => {
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

    // 3. Sort by most recent
    return Array.from(displayedCasesMap.values()).sort((a, b) => b.timestamp - a.timestamp);
  }, [cases, liveActions]);
};
