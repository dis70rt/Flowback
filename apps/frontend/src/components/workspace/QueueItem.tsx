import React from 'react';
import type { UnifiedQueueItem } from '../../hooks/useUnifiedQueue';

interface QueueItemProps {
  item: UnifiedQueueItem;
  isSelected: boolean;
  onClick: () => void;
}

export const QueueItem: React.FC<QueueItemProps> = ({ item, isSelected, onClick }) => {
  const riskAmount = item.amount_at_risk ? (item.amount_at_risk / 100).toLocaleString('en-IN') : '...';
  const initials = item.customer_id ? item.customer_id.slice(0, 2).toUpperCase() : item.case_id.slice(0, 2).toUpperCase();
  
  const isWhatsapp = item.channel === 'whatsapp' || item.channel === 'send_whatsapp';
  const isSms = item.channel === 'sms' || item.channel === 'send_sms';
  const isEmail = item.channel === 'email' || item.channel === 'send_email';
  const isCall = item.channel === 'call' || item.channel === 'send_call';

  let displayTag = item.channel || 'Pending';
  let statusColor = { dot: 'bg-rose-400', glow: 'shadow-[0_0_8px_rgba(251,113,133,0.5)]', bar: 'bg-rose-400', text: 'text-rose-300', badge: 'bg-rose-400/10 border-rose-400/20' };

  if (isWhatsapp) {
    displayTag = 'WhatsApp';
    statusColor = { dot: 'bg-emerald-400', glow: 'shadow-[0_0_8px_rgba(52,211,153,0.5)]', bar: 'bg-emerald-400', text: 'text-emerald-300', badge: 'bg-emerald-400/10 border-emerald-400/20' };
  } else if (isSms) {
    displayTag = 'SMS';
    statusColor = { dot: 'bg-[#0a84ff]', glow: 'shadow-[0_0_8px_rgba(10,132,255,0.5)]', bar: 'bg-[#0a84ff]', text: 'text-[#0a84ff]', badge: 'bg-[#0a84ff]/10 border-[#0a84ff]/20' };
  } else if (isEmail) {
    displayTag = 'Email';
    statusColor = { dot: 'bg-violet-400', glow: 'shadow-[0_0_8px_rgba(139,92,246,0.5)]', bar: 'bg-violet-400', text: 'text-violet-300', badge: 'bg-violet-400/10 border-violet-400/20' };
  } else if (isCall) {
    displayTag = 'Call';
    statusColor = { dot: 'bg-amber-400', glow: 'shadow-[0_0_8px_rgba(251,191,36,0.5)]', bar: 'bg-amber-400', text: 'text-amber-300', badge: 'bg-amber-400/10 border-amber-400/20' };
  }

  return (
    <div
      onClick={onClick}
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
            {new Date(item.timestamp).toLocaleTimeString('en-IN', { hour: 'numeric', minute: '2-digit' })}
          </span>
        </div>

        {/* Row 2: case ID */}
        <div className="flex items-center mb-1.5 mt-1">
          <span className="text-[10px] text-slate-500 font-medium mr-1.5">Case ID</span>
          <code className="text-[11px] font-medium text-slate-400 bg-slate-950/60 border border-slate-800/80 px-1.5 py-[2px] rounded-md font-mono shadow-sm">
            {item.case_id.slice(0, 12)}
          </code>
        </div>

        {/* Row 3: amount at risk */}
        <div className="flex items-center justify-between">
          <span className="text-[10.5px] text-slate-500">At Risk</span>
          <span className="text-[12px] font-bold text-rose-300">₹{riskAmount}</span>
        </div>
      </div>
    </div>
  );
};
