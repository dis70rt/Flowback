import { Camera, ChevronLeft, ChevronRight, Mic } from 'lucide-react';
import { formatSmsText, parseDraftBody } from './utils';

export const SmsMockup = ({ customer, pendingAction }: any) => {
  return (
    <>
      <div className="h-8 w-full bg-transparent flex justify-between items-end px-4 pb-1 text-[10.5px] font-bold text-white shrink-0 z-50">
        <span>9:41</span>
        <div className="flex items-center gap-1">
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><rect x="1" y="14" width="3" height="7" rx="0.5"/><rect x="6" y="10" width="3" height="11" rx="0.5"/><rect x="11" y="6" width="3" height="15" rx="0.5"/><rect x="16" y="2" width="3" height="19" rx="0.5" opacity="0.3"/></svg>
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><path d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2C16.93 2.93 7.08 2.93 1 9zm8 8l3 3 3-3c-1.65-1.66-4.34-1.66-6 0zm-4-4l2 2c2.76-2.76 7.24-2.76 10 0l2-2C15.14 9.14 8.87 9.14 5 13z"/></svg>
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none"><rect x="1.5" y="7" width="17" height="10" rx="2" stroke="white" strokeWidth="1.5"/><rect x="3" y="8.5" width="13" height="7" rx="1" fill="white"/><path d="M20 10.5v3" stroke="white" strokeWidth="1.5" strokeLinecap="round"/></svg>
        </div>
      </div>

      <div className="bg-[#1c1c1e]/90 backdrop-blur-md px-3 py-2 flex items-center gap-2.5 shrink-0 border-b border-white/5 z-40">
        <ChevronLeft className="w-6 h-6 text-[#0a84ff] cursor-pointer shrink-0 -ml-2" />
        <div className="flex flex-col items-center flex-1 -ml-4">
          <div className="w-8 h-8 rounded-full bg-slate-600 flex items-center justify-center overflow-hidden shrink-0">
             <div className="w-full h-full bg-[#6b7c85] flex items-center justify-center text-white text-[12px] font-bold">
                {customer?.name?.String?.slice(0, 2).toUpperCase() || 'CU'}
             </div>
          </div>
          <div className="flex items-center gap-1 mt-0.5">
            <span className="text-white text-[11px] font-semibold leading-tight truncate">{customer?.name?.String || 'Customer'}</span>
            <ChevronRight className="w-3 h-3 text-[#0a84ff]/80" />
          </div>
        </div>
      </div>

      <div className="flex-1 flex flex-col overflow-y-auto px-4 pt-3 pb-1 gap-1 bg-black">
        <div className="text-[#8e8e93] text-[10px] font-semibold text-center mb-4 mt-2">
          iMessage<br/>
          <span className="font-normal text-[9px]">Today 9:41 AM</span>
        </div>
        
        <div className="self-end max-w-[75%] flex flex-col items-end">
          <div className="relative text-white px-3.5 py-2 rounded-[18px] rounded-br-sm shadow-sm" style={{ backgroundColor: '#0a84ff' }}>
            <p className="text-[14.5px] leading-[1.35] whitespace-pre-wrap break-words">
              {formatSmsText(parseDraftBody(pendingAction?.draft_body))}
            </p>
            <svg className="absolute bottom-0 -right-[5px] text-[#0a84ff]" width="14" height="18" viewBox="0 0 14 18" fill="none">
               <path d="M0 18C4 18 8 16 10 11C10 14 12 17 14 18H0Z" fill="currentColor"/>
            </svg>
          </div>
          <div className="text-[#8e8e93] text-[10px] font-medium mt-1 pr-1 mr-1">Delivered</div>
        </div>
      </div>

      <div className="bg-[#1c1c1e] px-3 pt-2 pb-4 flex items-center gap-3 shrink-0">
        <Camera className="w-[22px] h-[22px] text-[#8e8e93] shrink-0" strokeWidth={1.5} />
        <div className="bg-black rounded-full flex-1 min-h-[32px] flex items-center px-3 border border-white/10">
          <span className="flex-1 text-[#8e8e93] text-[13px] select-none cursor-text">iMessage</span>
          <Mic className="w-[18px] h-[18px] text-[#8e8e93] shrink-0 ml-auto" strokeWidth={1.5} />
        </div>
      </div>
      
      <div className="absolute bottom-1.5 left-1/2 -translate-x-1/2 w-28 h-1 bg-white/40 rounded-full z-50" />
    </>
  );
};
