import { useState } from 'react';
import { useCases, useCaseDetails, useCustomer, useApproveAction, useRejectAction } from '../hooks/useApi';
import { SmsAppMockup, WhatsAppAppMockup } from '@/components/mockups';
import { useLiveActions } from '../hooks/useLiveActions';
import { useUnifiedQueue } from '../hooks/useUnifiedQueue';
import { QueueItem } from '../components/workspace/QueueItem';
import { motion, AnimatePresence } from 'framer-motion';

import { Badge } from '@/components/ui/badge';

import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Check, Edit2, Phone, Mail, MapPin, ChevronDown, ChevronUp } from 'lucide-react';

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
              <AnimatePresence initial={false}>
                {unifiedQueue.map((item) => (
                  <motion.div
                    key={item.case_id}
                    initial={{ opacity: 0, y: -20, scale: 0.95 }}
                    animate={{ opacity: 1, y: 0, scale: 1 }}
                    exit={{ opacity: 0, scale: 0.95 }}
                    transition={{ duration: 0.25, ease: "easeOut" }}
                    layout
                  >
                    <QueueItem 
                      item={item} 
                      isSelected={selectedCaseId === item.case_id} 
                      onClick={() => setSelectedCaseId(item.case_id)} 
                    />
                  </motion.div>
                ))}
              </AnimatePresence>
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

  const [isAiReasoningOpen, setIsAiReasoningOpen] = useState(false);
  const pendingAction = caseDetails?.actions?.find(a => a.status === 'PENDING' || a.status === 'PENDING_APPROVAL');

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

            <div className="relative w-full max-w-[380px] flex-1 min-h-[500px] my-4">
              <div
                className="absolute rounded-l-[3px]"
                style={{
                  left: -13,
                  top: 130,
                  width: 5,
                  height: 36,
                  background: "linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)",
                  boxShadow: "-2px 0 4px rgba(0,0,0,0.6), inset 1px 0 1px rgba(255,255,255,0.07)",
                  borderRadius: "3px 0 0 3px",
                }}
              />
              <div
                className="absolute"
                style={{
                  left: -13,
                  top: 178,
                  width: 5,
                  height: 36,
                  background: "linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)",
                  boxShadow: "-2px 0 4px rgba(0,0,0,0.6), inset 1px 0 1px rgba(255,255,255,0.07)",
                  borderRadius: "3px 0 0 3px",
                }}
              />
              <div
                className="absolute"
                style={{
                  right: -13,
                  top: 160,
                  width: 5,
                  height: 52,
                  background: "linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)",
                  boxShadow: "2px 0 4px rgba(0,0,0,0.6), inset -1px 0 1px rgba(255,255,255,0.07)",
                  borderRadius: "0 3px 3px 0",
                }}
              />
              <div
                className="absolute inset-0 rounded-[2.6rem] overflow-hidden"
                style={{
                  background: "linear-gradient(145deg, #2c2c2e 0%, #1c1c1e 35%, #111 65%, #0a0a0a 100%)",
                  boxShadow: "0 0 0 1.5px #444, 0 0 0 2.5px #1a1a1a, 0 25px 70px rgba(0,0,0,0.85), 0 8px 24px rgba(0,0,0,0.6), inset 0 1px 0 rgba(255,255,255,0.12), inset 0 -1px 0 rgba(0,0,0,0.5)",
                }}
              >
                <div
                  className="absolute top-0 left-6 right-6 h-[1.5px] rounded-full"
                  style={{ background: "linear-gradient(90deg, transparent, rgba(255,255,255,0.18), transparent)" }}
                />
                <div
                  className="absolute flex flex-col overflow-hidden"
                  style={{
                    top: 14,
                    left: 10,
                    right: 10,
                    bottom: 14,
                    borderRadius: "2rem",
                    background: pendingAction?.channel?.String === "send_sms" ? "#000000" : "#111b21",
                    boxShadow: "inset 0 0 0 1px rgba(0,0,0,0.8), 0 0 0 1px rgba(255,255,255,0.04)",
                  }}
                >
                  <div className="relative">
                    <div
                      className="absolute left-1/2 -translate-x-1/2 flex items-center justify-between"
                      style={{
                        top: 8,
                        width: 92,
                        height: 26,
                        padding: "0 8px",
                        zIndex: 30,
                        background: "#0a0a0a",
                        borderRadius: 20,
                        boxShadow: "inset 0 1px 3px rgba(0,0,0,0.9), 0 0 0 1px rgba(255,255,255,0.06)",
                      }}
                    >
                      <div style={{ width: 6, height: 6, borderRadius: "50%", background: "#111", boxShadow: "0 0 0 1px #1e1e1e" }} />
                      <div style={{ width: 12, height: 12, borderRadius: "50%", background: "radial-gradient(circle at 35% 35%, #1a2a3a, #060c14)", boxShadow: "0 0 0 1.5px #1a2535, inset 0 0 3px rgba(100,160,255,0.15)" }} />
                    </div>
                  </div>
                  {pendingAction?.channel?.String === "send_sms" ? (
                    <SmsAppMockup customer={customer} pendingAction={pendingAction} />
                  ) : (
                    <WhatsAppAppMockup customer={customer} pendingAction={pendingAction} />
                  )}
                  <div className="absolute inset-0 pointer-events-none rounded-[2rem]" style={{ background: "linear-gradient(135deg, rgba(255,255,255,0.04) 0%, transparent 45%)" }} />
                </div>
                <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-[4px]">
                  {Array.from({ length: 9 }).map((_, i) => (
                    <div key={i} style={{ width: 3, height: 3, borderRadius: "50%", background: "rgba(255,255,255,0.12)", boxShadow: "inset 0 1px 1px rgba(0,0,0,0.8)" }} />
                  ))}
                </div>
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
          {(caseDetails?.case?.status === 'PENDING_APPROVAL' || caseDetails?.case?.status === 'PENDING' || pendingAction) && (
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
              
              <div className="flex flex-col gap-2 mb-2">
                <Button 
                  size="default"
                  className="w-full h-10 px-4 bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm transition-all shadow-emerald-900/20"
                  onClick={() => approveAction.mutate(pendingAction.id)}
                  disabled={approveAction.isPending}
                >
                  <Check className="mr-2 w-4 h-4" /> Approve & Send
                </Button>
                
                <div className="flex w-full gap-2">
                  <Button 
                    size="default"
                    className="flex-1 h-10 px-4 shadow-sm transition-all bg-rose-600 hover:bg-rose-700 text-white"
                    onClick={() => rejectAction.mutate(pendingAction.id)}
                    disabled={rejectAction.isPending}
                  >
                    Reject
                  </Button>
                  <Button 
                    size="default"
                    className="shrink-0 h-10 px-4 shadow-sm transition-all bg-slate-800 hover:bg-slate-700 text-slate-200"
                    title="Edit Draft"
                  >
                    <Edit2 className="w-4 h-4 mr-2" /> Edit
                  </Button>
                </div>
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

              <div className="grid grid-cols-2 gap-4 pt-4 mt-3 border-t border-slate-800/60">
                <div>
                  <div className="text-[10px] uppercase tracking-wider text-slate-500">Payments</div>
                  <div className="font-medium text-slate-300 text-sm mt-0.5">
                    {customer?.successful_payments || 0} <span className="text-slate-600 text-[11px]">/ {customer?.total_payments || 0}</span>
                  </div>
                </div>
                <div>
                  <div className="text-[10px] uppercase tracking-wider text-slate-500">Tenure</div>
                  <div className="font-medium text-slate-300 text-sm mt-0.5 capitalize">{customer?.tenure?.String || 'N/A'}</div>
                </div>
                <div>
                  <div className="text-[10px] uppercase tracking-wider text-slate-500">Pref. Channel</div>
                  <div className="font-medium text-slate-300 text-sm mt-0.5 capitalize">{customer?.preferred_channel?.String?.replace('send_', '') || 'N/A'}</div>
                </div>
                <div>
                  <div className="text-[10px] uppercase tracking-wider text-slate-500">Reliability</div>
                  <div className="font-medium text-slate-300 text-sm mt-0.5 flex items-baseline gap-1">
                    {customer?.reliability_score || 0} <span className="text-slate-600 text-[10px]">/ 100</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Case Summary (Restored from Previous UI) */}
          <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 backdrop-blur-sm">
            <div className="text-[10px] text-slate-500 mb-4 uppercase font-semibold tracking-widest">Case Metrics</div>
            <div className="flex flex-col gap-3 text-[13px]">
              <div className="flex justify-between py-1.5 border-b border-slate-800/40">
                <span className="text-slate-500">Subscription ID</span>
                <span className="text-slate-300">{caseDetails?.case?.subscription_id}</span>
              </div>
              <div className="flex justify-between py-1.5 border-b border-slate-800/40">
                <span className="text-slate-500">Amount at Risk</span>
                <span className="text-rose-300 font-medium">₹{((caseDetails?.case?.amount_at_risk || 0) / 100).toLocaleString('en-IN')}</span>
              </div>
              <div className="flex justify-between py-1.5 border-b border-slate-800/40">
                <span className="text-slate-500">Failed Payments</span>
                <span className="text-rose-300 font-medium">{customer?.failed_payments || 0}</span>
              </div>
              <div className="flex justify-between py-1.5 border-b border-slate-800/40">
                <span className="text-slate-500">Status</span>
                <span className="text-slate-300">{caseDetails?.case?.status}</span>
              </div>
              <div className="flex justify-between py-1.5">
                <span className="text-slate-500">Created</span>
                <span className="text-slate-300">{caseDetails?.case?.created_at ? new Date(caseDetails.case.created_at).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : 'N/A'}</span>
              </div>
            </div>

            {/* AI Reasoning Expandable Section */}
            {pendingAction?.ai_reasoning?.Valid && (
              <div className="mt-4 pt-4 border-t border-slate-800/60">
                <button 
                  onClick={() => setIsAiReasoningOpen(!isAiReasoningOpen)}
                  className="flex items-center justify-between w-full text-left focus:outline-none"
                >
                  <span className="text-[10px] text-indigo-400 uppercase font-bold tracking-widest flex items-center gap-1.5">
                    <div className="w-1.5 h-1.5 rounded-full bg-indigo-400" /> AI Reasoning
                  </span>
                  {isAiReasoningOpen ? <ChevronUp className="w-3.5 h-3.5 text-slate-500" /> : <ChevronDown className="w-3.5 h-3.5 text-slate-500" />}
                </button>
                
                {isAiReasoningOpen && (
                  <div className="mt-3 text-[11.5px] text-slate-300/90 leading-relaxed bg-black/20 p-3.5 rounded-lg border border-slate-800/60 shadow-inner">
                    {pendingAction.ai_reasoning.String.split(/(ROOT CAUSE:|CONTEXT VIA TOOLS:|DECISION & COMPLIANCE:|NEXT STEP:|ACTION:|TIMING:|AUDIT:)/).map((part, i) => {
                       if (/(ROOT CAUSE:|CONTEXT VIA TOOLS:|DECISION & COMPLIANCE:|NEXT STEP:|ACTION:|TIMING:|AUDIT:)/.test(part)) {
                         return <div key={i} className="font-bold text-indigo-300 mt-2.5 mb-1 first:mt-0 tracking-wide text-[10px]">{part}</div>
                       }
                       return <span key={i}>{part}</span>
                    })}
                  </div>
                )}
              </div>
            )}
          </div>

        </div>
      </div>
    </div>
  );
};
