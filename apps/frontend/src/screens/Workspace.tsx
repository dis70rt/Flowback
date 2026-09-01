import { useState } from 'react';
import { useCases, useCaseDetails, useCustomer, useApproveAction, useRejectAction } from '../hooks/useApi';
import { SmsAppMockup, WhatsAppAppMockup } from '@/components/mockups';
import { useLiveActions } from '../hooks/useLiveActions';
import { useUnifiedQueue } from '../hooks/useUnifiedQueue';
import { QueueItem } from '../components/workspace/QueueItem';

import { Badge } from '@/components/ui/badge';

import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Check, Edit2, Phone, Mail, MapPin } from 'lucide-react';

export const Workspace = () => {
  const { actions } = useLiveActions();
  const { data: casesData, isPending: isLoadingCases } = useCases(1, 100);
  const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);

  const unifiedQueue = useUnifiedQueue(casesData?.data, actions);

  return (
    <div className="flex h-full w-full bg-slate-950">
      {/* Column 1: Live Queue */}
      <div className="w-[320px] border-r border-slate-800/60 flex flex-col h-full bg-slate-950">
        {/* Header */}
        <div className="px-4 pt-5 pb-3 border-b border-slate-800/60">
          <div className="flex items-center justify-between mb-1">
            <div className="flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-rose-400 shadow-[0_0_8px_rgba(251,113,133,0.7)] animate-pulse" />
              <h2 className="font-semibold text-[15px] text-white tracking-tight">Live Queue</h2>
            </div>
            <div className="bg-indigo-500/15 border border-indigo-500/25 text-indigo-300 text-[11px] font-bold px-2 py-0.5 rounded-full">
              {unifiedQueue.length} active
            </div>
          </div>
          <p className="text-[11px] text-slate-500">Payment recovery cases needing action</p>
        </div>

        {/* Cards */}
        <div className="flex-1 overflow-y-auto py-3 px-3">
          <div className="flex flex-col gap-2">
            {isLoadingCases ? (
              Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-[88px] w-full bg-slate-800/60 rounded-xl" />
              ))
            ) : (
              unifiedQueue.map((item) => (
                <QueueItem 
                  key={item.case_id}
                  item={item} 
                  isSelected={selectedCaseId === item.case_id} 
                  onClick={() => setSelectedCaseId(item.case_id)} 
                />
              ))
            )}
          </div>
        </div>
      </div>

      {/* Columns 2 & 3: Selected Case Details */}
      {selectedCaseId ? (
        <SelectedCaseView caseId={selectedCaseId} />
      ) : (
        <div className="flex-1 flex items-center justify-center text-slate-500">
          Select a case from the queue to view details
        </div>
      )}
    </div>
  );
};

const SelectedCaseView = ({ caseId }: { caseId: string }) => {
  const { data: caseDetails, isPending } = useCaseDetails(caseId);
  const { data: customer } = useCustomer(caseDetails?.case?.customer_id);
  
  const approveAction = useApproveAction();
  const rejectAction = useRejectAction();

  const pendingAction = caseDetails?.actions?.find(a => a.status === 'PENDING_APPROVAL');

  if (isPending) {
    return (
      <div className="flex-1 flex">
        {/* Center Column Skeleton */}
        <div className="flex-1 border-r border-slate-800 p-6 bg-slate-950 flex flex-col items-center justify-center overflow-hidden">
          <Skeleton className="w-[380px] h-[760px] rounded-[3.5rem] bg-slate-900/40 border-[14px] border-slate-900/60 shadow-2xl" />
        </div>

        {/* Right Column Skeleton */}
        <div className="w-[380px] p-6 flex flex-col gap-6 bg-slate-950">
          <div className="flex items-center justify-between mb-2">
            <Skeleton className="h-6 w-32 bg-slate-800/80" />
            <Skeleton className="h-6 w-16 bg-slate-800/80" />
          </div>
          <Skeleton className="h-[200px] w-full bg-slate-900/60 rounded-xl" />
          <Skeleton className="h-[300px] w-full bg-slate-900/60 rounded-xl" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex">
      {/* Column 2: Action Preview (Now in Center) */}
      <div className="flex-1 border-r border-slate-800 p-6 bg-slate-950 flex flex-col items-center overflow-hidden">
        {pendingAction ? (
          <>
            <div className="w-full max-w-[380px] flex flex-col mb-4 gap-1 shrink-0">
              <div className="text-[10px] uppercase tracking-widest text-slate-500 font-semibold">{pendingAction?.channel?.String === 'send_sms' ? 'SMS / iMessage' : 'WhatsApp'}</div>
              <div className="flex justify-between items-center">
                <span className="text-sm font-medium text-slate-200 tracking-wide">AI RECOVERY DRAFT</span>
                <span className="text-[10px] text-emerald-200/70 font-medium tracking-wider flex items-center gap-1.5 bg-emerald-900/20 px-2 py-0.5 rounded-full border border-emerald-800/30">
                  <Check className="w-3 h-3" /> READY FOR REVIEW
                </span>
              </div>
            </div>

            <div className="flex-1 overflow-y-auto w-full flex justify-center pb-8 hide-scrollbar">
              <div className="relative w-[380px] shrink-0">
                
                {pendingAction?.channel?.String === 'send_sms' ? (
                  <SmsAppMockup customer={customer} pendingAction={pendingAction} />
                ) : (
                  <WhatsAppAppMockup customer={customer} pendingAction={pendingAction} />
                )}

              </div>
            </div>
          </>
        ) : (
          <div className="h-full flex flex-col items-center justify-center text-slate-500 gap-3">
            <div className="w-12 h-12 rounded-full bg-slate-900 border border-slate-800 flex items-center justify-center">
              <Check className="w-6 h-6 text-emerald-500/50" />
            </div>
            <p>No pending actions for this case</p>
          </div>
        )}
      </div>

      {/* Column 3: Context Panel (Now on Right) */}
      <div className="w-[380px] bg-slate-950 flex flex-col h-full">
        {/* Header */}
        <div className="p-6 border-b border-slate-800/60 flex items-center justify-between shrink-0">
          <h2 className="text-sm font-semibold text-white tracking-wide">Review Context</h2>
          {caseDetails?.case?.status === 'PENDING_APPROVAL' && (
            <Badge variant="outline" className="bg-amber-400/10 text-amber-300 border-amber-400/20 px-2 py-0">
              Needs Review
            </Badge>
          )}
        </div>

        {/* Scrollable Content */}
        <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-8">
          
          {/* Action Required Block */}
          {pendingAction && (
            <div className="bg-indigo-950/20 border border-indigo-500/20 rounded-xl p-5 shadow-lg shadow-indigo-500/5">
              <h3 className="text-xs uppercase font-bold text-indigo-300 tracking-wider mb-4 flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-indigo-400 animate-pulse" />
                Action Required
              </h3>
              
              <div className="flex gap-2 mb-4">
                <Button 
                  className="flex-1 bg-emerald-500 hover:bg-emerald-600 text-white shadow-sm"
                  onClick={() => approveAction.mutate(pendingAction.id)}
                  disabled={approveAction.isPending}
                >
                  <Check className="w-4 h-4 mr-2" /> Approve & Send
                </Button>
                <Button 
                  variant="outline"
                  className="px-3 border-slate-700 hover:bg-slate-800 text-slate-300"
                  onClick={() => rejectAction.mutate(pendingAction.id)}
                  disabled={rejectAction.isPending}
                >
                  <Edit2 className="w-4 h-4" />
                </Button>
              </div>

              <div className="text-[11px] text-slate-400 leading-relaxed border-t border-indigo-500/10 pt-4">
                <span className="font-semibold text-indigo-300/80 mr-1">AI Reasoning:</span>
                {pendingAction.ai_reasoning?.String || 'Determined best time and channel to reach out based on previous successful recoveries.'}
              </div>
            </div>
          )}

          {/* Customer Context */}
          <div className="space-y-4">
            <h3 className="text-xs uppercase font-bold text-slate-500 tracking-wider">Customer Profile</h3>
            <div className="bg-slate-900/40 border border-slate-800/60 rounded-xl p-4 space-y-3">
              <div className="flex justify-between items-start">
                <span className="text-sm font-medium text-slate-200">{customer?.name?.String || 'Unknown'}</span>
                {customer?.value_tier?.Valid && (
                  <Badge variant="outline" className="bg-slate-800/50 text-slate-400 border-slate-700/50">
                    {customer.value_tier.String} Tier
                  </Badge>
                )}
              </div>

              <div className="space-y-2 mt-4">
                <div className="flex items-center text-[11px] text-slate-400">
                  <Mail className="w-3.5 h-3.5 mr-2 text-slate-500" />
                  {customer?.email?.String || 'No email provided'}
                </div>
                <div className="flex items-center text-[11px] text-slate-400">
                  <Phone className="w-3.5 h-3.5 mr-2 text-slate-500" />
                  {customer?.phone?.String || 'No phone provided'}
                </div>
                <div className="flex items-center text-[11px] text-slate-400">
                  <MapPin className="w-3.5 h-3.5 mr-2 text-slate-500" />
                  {[customer?.city?.String, customer?.state?.String].filter(Boolean).join(', ') || 'Unknown Location'}
                </div>
              </div>

              <div className="border-t border-slate-800/60 pt-3 mt-3 flex items-center justify-between">
                <div>
                  <div className="text-[10px] text-slate-500 font-medium">Reliability Score</div>
                  <div className="text-sm font-bold text-emerald-400">{customer?.reliability_score || 0}/100</div>
                </div>
                <div className="text-right">
                  <div className="text-[10px] text-slate-500 font-medium">Success Rate</div>
                  <div className="text-sm font-bold text-slate-200">
                    {customer?.total_payments ? Math.round((customer.successful_payments / customer.total_payments) * 100) : 0}%
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
};
