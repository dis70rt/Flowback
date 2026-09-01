import React, { createContext, useContext, useEffect, useState } from 'react';
import { toast } from 'sonner';
import { useQueryClient } from '@tanstack/react-query';

export interface LiveAction {
  event: string;
  case_id: string;
  customer_id: string;
  channel: string;
  status: string;
  audio_url?: string;
  timestamp: number;
}

interface LiveActionsContextType {
  actions: LiveAction[];
  removeAction: (caseId: string) => void;
}

const LiveActionsContext = createContext<LiveActionsContextType | undefined>(undefined);

export const LiveActionsProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [actions, setActions] = useState<LiveAction[]>([]);
  const queryClient = useQueryClient();

  useEffect(() => {
    const sse = new EventSource('/api/stream');
    
    sse.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.event === 'draft_ready') {
          setActions(prev => {
            if (prev.some(a => a.case_id === data.case_id)) return prev;
            return [{ ...data, timestamp: Date.now() }, ...prev];
          });
          toast('New Draft Ready', {
            description: `Case ${data.case_id?.slice(0, 8)} needs review for ${data.channel}`,
          });
        }
        
        queryClient.invalidateQueries({ queryKey: ['metrics'] });
        queryClient.invalidateQueries({ queryKey: ['cases'] });
      } catch (e) {
        console.error('SSE parse error:', e);
      }
    };

    return () => sse.close();
  }, [queryClient]);

  const removeAction = (caseId: string) => {
    setActions(prev => prev.filter(a => a.case_id !== caseId));
  };

  return (
    <LiveActionsContext.Provider value={{ actions, removeAction }}>
      {children}
    </LiveActionsContext.Provider>
  );
};

export const useLiveActions = () => {
  const context = useContext(LiveActionsContext);
  if (!context) {
    throw new Error('useLiveActions must be used within a LiveActionsProvider');
  }
  return context;
};
