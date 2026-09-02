import { ArrowLeft, Archive, Trash2, Mail, MoreVertical, Reply, CornerUpRight, Star } from 'lucide-react';
import { formatEmailText } from './utils';

export const EmailMockup = ({ customer, pendingAction }: any) => {
  const subject = pendingAction?.draft_payload?.RawMessage?.subject || 'Important update regarding your account';
  const body = pendingAction?.draft_payload?.RawMessage?.body || '';

  return (
    <div className="flex flex-col h-full w-full bg-[#121212] font-sans relative">
      {/* Status Bar */}
      <div className="h-8 w-full bg-transparent flex justify-between items-end px-4 pb-1 text-[10.5px] font-bold text-white shrink-0 z-50">
        <span>9:41</span>
        <div className="flex items-center gap-1">
          {/* Signal */}
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><rect x="1" y="14" width="3" height="7" rx="0.5"/><rect x="6" y="10" width="3" height="11" rx="0.5"/><rect x="11" y="6" width="3" height="15" rx="0.5"/><rect x="16" y="2" width="3" height="19" rx="0.5" opacity="0.3"/></svg>
          {/* Wifi */}
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><path d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2C16.93 2.93 7.08 2.93 1 9zm8 8l3 3 3-3c-1.65-1.66-4.34-1.66-6 0zm-4-4l2 2c2.76-2.76 7.24-2.76 10 0l2-2C15.14 9.14 8.87 9.14 5 13z"/></svg>
          {/* Battery */}
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none"><rect x="1.5" y="7" width="17" height="10" rx="2" stroke="white" strokeWidth="1.5"/><rect x="3" y="8.5" width="13" height="7" rx="1" fill="white"/><path d="M20 10.5v3" stroke="white" strokeWidth="1.5" strokeLinecap="round"/></svg>
        </div>
      </div>

      {/* App Bar */}
      <div className="flex items-center justify-between px-4 py-3 shrink-0 text-[#e3e3e3] mt-2 z-40">
        <ArrowLeft className="w-6 h-6 text-[#e3e3e3]" />
        <div className="flex items-center gap-6">
          <Archive className="w-5 h-5 text-[#e3e3e3]" />
          <Trash2 className="w-5 h-5 text-[#e3e3e3]" />
          <Mail className="w-5 h-5 text-[#e3e3e3]" />
          <MoreVertical className="w-5 h-5 text-[#e3e3e3]" />
        </div>
      </div>

      {/* Scrollable Content */}
      <div className="flex-1 overflow-y-auto px-4 pb-20 hide-scrollbar flex flex-col z-30">
        {/* Subject */}
        <div className="flex justify-between items-start mt-4 mb-6">
          <h1 className="text-[22px] font-normal text-[#e3e3e3] leading-snug pr-4">{subject}</h1>
          <Star className="w-6 h-6 text-[#c4c7c5] shrink-0 mt-0.5" />
        </div>

        {/* Sender Info */}
        <div className="flex gap-3 mb-6">
          <div className="w-11 h-11 rounded-full bg-[#20516d] flex items-center justify-center text-white text-[18px] font-medium shrink-0">
            F
          </div>
          <div className="flex flex-col flex-1 justify-center">
            <div className="flex justify-between items-baseline">
              <span className="text-[#e3e3e3] font-medium text-[15px]">FlowBack Support</span>
              <span className="text-[#c4c7c5] text-[12px]">Just now</span>
            </div>
            <div className="flex items-center gap-1 text-[#c4c7c5] text-[12px] mt-0.5">
              <span>to {customer?.name?.String?.split(" ")[0] || "me"}</span>
              <svg className="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24"><path d="M7 10l5 5 5-5z"/></svg>
            </div>
          </div>
        </div>

        {/* Email Body */}
        <div className="text-[#e3e3e3] text-[15px] leading-relaxed font-sans whitespace-pre-wrap">
          {formatEmailText(body)}
        </div>

        {/* Action Pills */}
        <div className="flex gap-2 mt-10">
          <div className="flex items-center gap-1.5 px-3 py-1.5 bg-[#2d2d2d] rounded-full text-[#e3e3e3] text-[12.5px] font-medium border border-[#3d3d3d]">
            <Reply className="w-3.5 h-3.5 text-[#c4c7c5]" /> Reply
          </div>
          <div className="flex items-center gap-1.5 px-3 py-1.5 bg-[#2d2d2d] rounded-full text-[#e3e3e3] text-[12.5px] font-medium border border-[#3d3d3d]">
            <Reply className="w-3.5 h-3.5 text-[#c4c7c5] scale-x-[-1]" /> Reply all
          </div>
          <div className="flex items-center gap-1.5 px-3 py-1.5 bg-[#2d2d2d] rounded-full text-[#e3e3e3] text-[12.5px] font-medium border border-[#3d3d3d]">
            <CornerUpRight className="w-3.5 h-3.5 text-[#c4c7c5]" /> Forward
          </div>
        </div>
      </div>

      {/* Home Indicator */}
      <div className="absolute bottom-1.5 left-1/2 -translate-x-1/2 w-28 h-1 bg-white/30 rounded-full z-50" />
    </div>
  );
};
