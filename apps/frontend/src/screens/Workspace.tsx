import { useState } from 'react';
import { useCases, useCaseDetails, useCustomer, useApproveAction, useRejectAction } from '../hooks/useApi';
import { SmsAppMockup, WhatsAppAppMockup } from '@/components/mockups';

import { Badge } from '@/components/ui/badge';

import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Check, Edit2, Phone, Mail, MapPin } from 'lucide-react';

export const Workspace = () => {
  const { data: casesData, isPending: isLoadingCases } = useCases(1, 20);
  const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);

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
              {casesData?.data?.length || 0} active
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
              casesData?.data?.map((item) => {
                const isSelected = selectedCaseId === item.id;
                const riskAmount = (item.amount_at_risk / 100).toLocaleString('en-IN');
                const initials = item.id.slice(0, 2).toUpperCase();
                const isDetected = item.status === 'DETECTED' || item.status === 'send_whatsapp';
                const displayTag = isDetected ? 'WhatsApp' : item.status;
                const statusColor =
                  isDetected                      ? { dot: 'bg-emerald-400', glow: 'shadow-[0_0_8px_rgba(52,211,153,0.5)]', bar: 'bg-emerald-400', text: 'text-emerald-300', badge: 'bg-emerald-400/10 border-emerald-400/20' } :
                  item.status === 'In Progress'   ? { dot: 'bg-cyan-400',  glow: 'shadow-[0_0_8px_rgba(34,211,238,0.5)]',  bar: 'bg-cyan-400',  text: 'text-cyan-300',  badge: 'bg-cyan-400/10 border-cyan-400/20'  } :
                  item.status === 'Pending Review'? { dot: 'bg-amber-400', glow: 'shadow-[0_0_8px_rgba(251,191,36,0.5)]',  bar: 'bg-amber-400', text: 'text-amber-300', badge: 'bg-amber-400/10 border-amber-400/20' } :
                                                   { dot: 'bg-rose-400',glow: 'shadow-[0_0_8px_rgba(251,113,133,0.5)]', bar: 'bg-rose-400',text: 'text-rose-300',badge: 'bg-rose-400/10 border-rose-400/20'};

                return (
                  <div
                    key={item.id}
                    onClick={() => setSelectedCaseId(item.id)}
                    className={`
                      group relative flex gap-3 p-3 rounded-xl cursor-pointer
                      border transition-all duration-150 overflow-hidden
                      ${isSelected
                        ? 'bg-slate-800/80 border-slate-600/60 shadow-[0_0_0_1px_rgba(139,92,246,0.3)]'
                        : 'bg-slate-900/40 border-slate-800/60 hover:bg-slate-800/50 hover:border-slate-700/60'}
                    `}
                  >
                    {/* Status bar accent on left edge */}
                    {isSelected && (
                      <div className={`absolute left-0 top-3 bottom-3 w-[3px] rounded-full ${statusColor.bar} opacity-80`} />
                    )}

                    {/* Avatar */}
                    <div className="w-9 h-9 rounded-lg bg-slate-800 border border-slate-700/50 flex items-center justify-center shrink-0 text-[13px] font-bold text-slate-300 ml-2">
                      {initials}
                    </div>

                    {/* Content */}
                    <div className="flex-1 min-w-0">
                      {/* Row 1: status + date */}
                      <div className="flex items-center justify-between mb-1">
                        <div className={`flex items-center gap-1.5 text-[10.5px] font-semibold border rounded-full px-2 py-[1px] ${statusColor.text} ${statusColor.badge}`}>
                          <div className={`w-1.5 h-1.5 rounded-full ${statusColor.dot} ${statusColor.glow}`} />
                          {displayTag}
                        </div>
                        <span className="text-[10px] text-slate-600 font-medium">
                          {new Date(item.created_at).toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })}
                        </span>
                      </div>

                      {/* Row 2: case ID */}
                      <div className="flex items-center mb-1.5 mt-1">
                        <span className="text-[10px] text-slate-500 font-medium mr-1.5">Case ID</span>
                        <code className="text-[11px] font-medium text-slate-400 bg-slate-950/60 border border-slate-800/80 px-1.5 py-[2px] rounded-md font-mono shadow-sm">
                          {item.id.slice(0, 12)}
                        </code>
                      </div>

                      {/* Row 3: amount at risk */}
                      <div className="flex items-center justify-between">
                        <span className="text-[10.5px] text-slate-500">At Risk</span>
                        <span className="text-[12px] font-bold text-rose-300">₹{riskAmount}</span>
                      </div>
                    </div>

                    {/* Selected indicator chevron */}
                    {isSelected && (
                      <div className="absolute right-2 top-1/2 -translate-y-1/2 w-1 h-4 rounded-full bg-violet-400/60" />
                    )}
                  </div>
                );
              })
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
  const { data: caseDetails, isPending: isLoadingCase } = useCaseDetails(caseId);
  const customerId = caseDetails?.case?.customer_id;
  const { data: customer, isPending: isLoadingCustomer } = useCustomer(customerId);

  const approveMutation = useApproveAction();
  const rejectMutation = useRejectAction();

  const pendingAction = caseDetails?.actions?.find(a => a.status === 'PENDING_APPROVAL' || a.status === 'draft' || a.status === 'PENDING');

  if (isLoadingCase) {
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
            <Skeleton className="h-6 w-24 bg-slate-800/80 rounded-full" />
          </div>
          
          {/* Action Card Skeleton */}
          <Skeleton className="h-[140px] w-full bg-slate-900/60 rounded-xl border border-slate-800/50" />

          {/* Customer Card Skeleton */}
          <Skeleton className="h-[210px] w-full bg-slate-900/40 rounded-xl border border-slate-800/40" />

          {/* Case Metrics Skeleton */}
          <Skeleton className="h-[220px] w-full bg-slate-900/40 rounded-xl border border-slate-800/40" />
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
                  <span className="w-1.5 h-1.5 bg-emerald-400/80 rounded-full animate-pulse shadow-[0_0_5px_rgba(52,211,153,0.6)]"></span> READY
                </span>
              </div>
            </div>
            
            {/* ── Realistic Phone Shell ─────────────────────────── */}
            {/* Outer wrapper — responsive to fill available height */}
            <div className="relative w-full max-w-[380px] flex-1 min-h-[500px]">

              {/* LEFT side buttons */}
              {/* Volume Up */}
              <div
                className="absolute rounded-l-[3px]"
                style={{
                  left: -13,
                  top: 130,
                  width: 5,
                  height: 36,
                  background: 'linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)',
                  boxShadow: '-2px 0 4px rgba(0,0,0,0.6), inset 1px 0 1px rgba(255,255,255,0.07)',
                  borderRadius: '3px 0 0 3px',
                }}
              />
              {/* Volume Down */}
              <div
                className="absolute"
                style={{
                  left: -13,
                  top: 178,
                  width: 5,
                  height: 36,
                  background: 'linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)',
                  boxShadow: '-2px 0 4px rgba(0,0,0,0.6), inset 1px 0 1px rgba(255,255,255,0.07)',
                  borderRadius: '3px 0 0 3px',
                }}
              />

              {/* RIGHT side — Power button */}
              <div
                className="absolute"
                style={{
                  right: -13,
                  top: 160,
                  width: 5,
                  height: 52,
                  background: 'linear-gradient(180deg, #3a3a3a 0%, #2a2a2a 40%, #222 100%)',
                  boxShadow: '2px 0 4px rgba(0,0,0,0.6), inset -1px 0 1px rgba(255,255,255,0.07)',
                  borderRadius: '0 3px 3px 0',
                }}
              />

              {/* Phone chassis — the body */}
              <div
                className="absolute inset-0 rounded-[2.6rem] overflow-hidden"
                style={{
                  background: 'linear-gradient(145deg, #2c2c2e 0%, #1c1c1e 35%, #111 65%, #0a0a0a 100%)',
                  boxShadow: `
                    0 0 0 1.5px #444,
                    0 0 0 2.5px #1a1a1a,
                    0 25px 70px rgba(0,0,0,0.85),
                    0 8px 24px rgba(0,0,0,0.6),
                    inset 0 1px 0 rgba(255,255,255,0.12),
                    inset 0 -1px 0 rgba(0,0,0,0.5)
                  `,
                }}
              >
                {/* Top glossy edge highlight */}
                <div
                  className="absolute top-0 left-6 right-6 h-[1.5px] rounded-full"
                  style={{ background: 'linear-gradient(90deg, transparent, rgba(255,255,255,0.18), transparent)' }}
                />

                {/* Screen inset — sits inside the chassis with a bezel */}
                <div
                  className="absolute flex flex-col overflow-hidden"
                  style={{
                    top: 14,
                    left: 10,
                    right: 10,
                    bottom: 14,
                    borderRadius: '2rem',
                    background: pendingAction?.channel?.String === 'send_sms' ? '#000000' : '#111b21',
                    boxShadow: 'inset 0 0 0 1px rgba(0,0,0,0.8), 0 0 0 1px rgba(255,255,255,0.04)',
                  }}
                >
                  {/* Camera punch-hole pill — top centre (Wider Dynamic Island style) */}
                  <div className="relative">
                    <div
                      className="absolute left-1/2 -translate-x-1/2 flex items-center justify-between"
                      style={{
                        top: 8,
                        width: 92,
                        height: 26,
                        padding: '0 8px',
                        zIndex: 30,
                        background: '#0a0a0a',
                        borderRadius: 20,
                        boxShadow: 'inset 0 1px 3px rgba(0,0,0,0.9), 0 0 0 1px rgba(255,255,255,0.06)',
                      }}
                    >
                      {/* Flash/sensor dot (Left) */}
                      <div
                        style={{
                          width: 6,
                          height: 6,
                          borderRadius: '50%',
                          background: '#111',
                          boxShadow: '0 0 0 1px #1e1e1e',
                        }}
                      />
                      {/* Camera lens (Right) */}
                      <div
                        style={{
                          width: 12,
                          height: 12,
                          borderRadius: '50%',
                          background: 'radial-gradient(circle at 35% 35%, #1a2a3a, #060c14)',
                          boxShadow: '0 0 0 1.5px #1a2535, inset 0 0 3px rgba(100,160,255,0.15)',
                        }}
                      />
                    </div>
                  </div>

                  {pendingAction?.channel?.String === 'send_sms' ? (
                    <SmsAppMockup customer={customer} pendingAction={pendingAction} />
                  ) : (
                    <WhatsAppAppMockup customer={customer} pendingAction={pendingAction} />
                  )}

                  {/* Subtle screen glare overlay */}
                  <div
                    className="absolute inset-0 pointer-events-none rounded-[2rem]"
                    style={{
                      background: 'linear-gradient(135deg, rgba(255,255,255,0.04) 0%, transparent 45%)',
                    }}
                  />
                </div>

                {/* Bottom speaker grille dots */}
                <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-[4px]">
                  {Array.from({ length: 9 }).map((_, i) => (
                    <div
                      key={i}
                      style={{
                        width: 3,
                        height: 3,
                        borderRadius: '50%',
                        background: 'rgba(255,255,255,0.12)',
                        boxShadow: 'inset 0 1px 1px rgba(0,0,0,0.8)',
                      }}
                    />
                  ))}
                </div>
              </div>
            </div>
          </>
        ) : (
           <div className="flex h-full items-center justify-center flex-col text-slate-500 text-center p-8">
             <div className="w-20 h-20 rounded-full bg-slate-900 flex items-center justify-center mb-6 border border-slate-800">
               <Check className="w-8 h-8 text-slate-600" />
             </div>
             <h3 className="text-slate-300 font-medium mb-2 text-lg">No Action Required</h3>
             <p className="text-sm">There are no pending draft actions for this case.</p>
           </div>
        )}
      </div>

      {/* Column 3: Customer & Case Summary (Now on the Right) */}
      <div className="w-[380px] p-6 flex flex-col gap-6 overflow-y-auto bg-slate-950">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-medium text-slate-200">Case Context</h2>
          <Badge className="bg-slate-800 text-slate-300 border-slate-700">{caseDetails?.case?.status}</Badge>
        </div>

        {/* Action Card (Moved here) */}
        {pendingAction && (
          <div className="bg-slate-900/80 rounded-xl border border-slate-800/80 p-5 flex flex-col gap-3 shadow-lg">
             <div className="text-[10px] text-slate-500 uppercase font-semibold tracking-widest mb-1">Required Action</div>
             <Button 
                className="w-full bg-emerald-600/20 text-emerald-300 border border-emerald-500/30 hover:bg-emerald-600/30 transition-colors shadow-[0_0_10px_rgba(52,211,153,0.1)] py-5"
                onClick={() => approveMutation.mutate(pendingAction.id)}
                disabled={approveMutation.isPending}
             >
                <Check className="w-5 h-5 mr-2" /> Approve & Send Message
             </Button>
             <Button 
                variant="outline" 
                className="w-full border-rose-500/30 text-rose-300 hover:bg-rose-500/10 hover:text-rose-200 bg-transparent transition-colors py-5"
                onClick={() => rejectMutation.mutate(pendingAction.id)}
                disabled={rejectMutation.isPending}
             >
                <Edit2 className="w-4 h-4 mr-2" /> Reject / Edit Draft
             </Button>
          </div>
        )}

        {/* Customer Card */}
        <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 backdrop-blur-sm">
          <div className="text-[10px] text-slate-500 mb-4 uppercase font-semibold tracking-widest">Customer Profile</div>
          {isLoadingCustomer ? (
            <Skeleton className="h-20 w-full bg-slate-800" />
          ) : (
            <>
              <div className="flex items-center justify-between mb-5">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-full bg-indigo-900/40 text-indigo-300 flex items-center justify-center text-lg font-bold border border-indigo-500/20">
                    {customer?.name?.String?.slice(0, 2).toUpperCase() || 'C'}
                  </div>
                  <div>
                    <div className="font-medium text-slate-200 text-[15px]">{customer?.name?.String || 'Unknown Customer'}</div>
                    <div className="text-[11px] text-slate-500 font-mono mt-0.5">ID: {customer?.id.slice(0,8)}...</div>
                  </div>
                </div>
                {customer?.value_tier?.Valid && (
                  <Badge variant="outline" className="text-indigo-300/80 border-indigo-500/20 bg-indigo-500/10 px-2 py-0.5 rounded-full text-[10px]">{customer.value_tier.String}</Badge>
                )}
              </div>

              {/* Contact Information */}
              <div className="flex flex-col gap-2 mb-5">
                {customer?.phone?.Valid && (
                  <div className="flex items-center gap-2 text-[12px] text-slate-300">
                    <Phone className="w-3.5 h-3.5 text-slate-500" />
                    {customer.phone.String}
                  </div>
                )}
                {customer?.email?.Valid && (
                  <div className="flex items-center gap-2 text-[12px] text-slate-300">
                    <Mail className="w-3.5 h-3.5 text-slate-500" />
                    {customer.email.String}
                  </div>
                )}
                {(customer?.city?.Valid || customer?.state?.Valid) && (
                  <div className="flex items-center gap-2 text-[12px] text-slate-300">
                    <MapPin className="w-3.5 h-3.5 text-slate-500" />
                    {[customer?.city?.String, customer?.state?.String].filter(Boolean).join(', ')}
                  </div>
                )}
              </div>

              {/* Stats Grid */}
              <div className="grid grid-cols-2 gap-4 pt-4 border-t border-slate-800/60">
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
            </>
          )}
        </div>

        {/* Case Summary */}
        <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 backdrop-blur-sm">
          <div className="text-[10px] text-slate-500 mb-4 uppercase font-semibold tracking-widest">Case Metrics</div>
          <div className="flex flex-col gap-3 text-[13px]">
            <div className="flex justify-between py-1.5 border-b border-slate-800/40">
              <span className="text-slate-500">Subscription ID</span>
              <span className="text-slate-300">{caseDetails?.case?.subscription_id}</span>
            </div>
            <div className="flex justify-between py-1.5 border-b border-slate-800/40">
              <span className="text-slate-500">Amount at Risk</span>
              <span className="text-rose-300 font-medium">₹{((caseDetails?.case?.amount_at_risk || 0) / 100).toLocaleString()}</span>
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
              <span className="text-slate-300">{new Date(caseDetails?.case?.created_at || '').toLocaleString()}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
