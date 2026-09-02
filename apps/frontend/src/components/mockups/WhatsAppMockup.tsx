import { ArrowLeft, MoreVertical, Phone, Video, Smile, Paperclip, Mic, Camera } from 'lucide-react';
import { formatWhatsAppText } from './utils';

export const WhatsAppMockup = ({ customer, pendingAction }: any) => {
  return (
    <>
      <div className="h-8 w-full bg-[#075e54] flex justify-between items-end px-4 pb-1 text-[10.5px] font-bold text-white/95 shrink-0 z-50">
        <span>9:41</span>
        <div className="flex items-center gap-1">
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><rect x="1" y="14" width="3" height="7" rx="0.5"/><rect x="6" y="10" width="3" height="11" rx="0.5"/><rect x="11" y="6" width="3" height="15" rx="0.5"/><rect x="16" y="2" width="3" height="19" rx="0.5" opacity="0.3"/></svg>
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><path d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2C16.93 2.93 7.08 2.93 1 9zm8 8l3 3 3-3c-1.65-1.66-4.34-1.66-6 0zm-4-4l2 2c2.76-2.76 7.24-2.76 10 0l2-2C15.14 9.14 8.87 9.14 5 13z"/></svg>
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none"><rect x="1.5" y="7" width="17" height="10" rx="2" stroke="white" strokeWidth="1.5"/><rect x="3" y="8.5" width="13" height="7" rx="1" fill="white"/><path d="M20 10.5v3" stroke="white" strokeWidth="1.5" strokeLinecap="round"/></svg>
        </div>
      </div>

      <div className="bg-[#075e54] px-3 py-2 flex items-center gap-2.5 shrink-0 z-40">
        <ArrowLeft className="w-5 h-5 text-white/90 cursor-pointer shrink-0" strokeWidth={2.5} />
        <div className="w-9 h-9 rounded-full bg-[#3b4a54] flex items-center justify-center overflow-hidden shrink-0 border border-white/10">
          <div className="w-full h-full bg-[#6b7c85] flex items-center justify-center text-white text-[13px] font-bold">
            {customer?.name?.String?.slice(0, 2).toUpperCase() || 'CU'}
          </div>
        </div>
        <div className="flex-1 min-w-0">
          <div className="text-white text-[15px] font-semibold leading-tight truncate">
            {customer?.name?.String || 'Customer'}
          </div>
          <div className="text-[#b2dfdb] text-[12px] leading-tight">online</div>
        </div>
        <div className="flex items-center gap-4 text-white/90">
          <Video className="w-[19px] h-[19px] cursor-pointer" strokeWidth={2} />
          <Phone className="w-[18px] h-[18px] cursor-pointer" strokeWidth={2} />
          <MoreVertical className="w-[18px] h-[18px] cursor-pointer" strokeWidth={2} />
        </div>
      </div>

      <div
        className="flex-1 flex flex-col overflow-y-auto px-2 pt-3 pb-1 gap-1"
        style={{
          backgroundImage: `url('/wa-wallpaper.png')`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
        }}
      >
        <div className="self-center bg-[#182229]/90 text-[#8696a0] text-[10.5px] font-semibold px-3 py-1 rounded-lg mb-2 shadow-sm backdrop-blur-sm">
          TODAY
        </div>
        <div className="self-center bg-[#182229]/90 text-[#ffd279] text-[10px] text-center px-3 py-2 rounded-lg mb-2 max-w-[85%] shadow-sm backdrop-blur-sm leading-snug">
          🔒 Messages to this business are now secured with end-to-end encryption.
        </div>
        <div className="self-end max-w-[80%] flex flex-col items-end">
          <div className="relative text-[#e9edef] px-3 pt-2 pb-1.5 rounded-[8px] rounded-tr-none shadow-md" style={{ backgroundColor: '#005c4b' }}>
            <svg className="absolute top-0 -right-[8px] text-[#005c4b]" width="8" height="13" viewBox="0 0 8 13" fill="none">
              <path fill="currentColor" d="M5.188 0H0v11.193l6.467-8.625C7.526 1.156 6.958 0 5.188 0z" />
            </svg>
            <p className="text-[13.5px] leading-[1.45] whitespace-pre-wrap break-words">
              {formatWhatsAppText(pendingAction?.draft_payload?.RawMessage?.body || "")}
            </p>
            <div className="flex items-center justify-end mt-0.5">
              <span className="text-[#8696a0] text-[10.5px]">9:41 AM</span>
              <svg viewBox="0 0 16 11" width="16" height="11" className="fill-[#53bdeb] ml-1">
                <path d="M11.071.653l-.542-.45a.394.394 0 0 0-.553.063L5.115 6.505l-1.95-1.836a.417.417 0 0 0-.59.008l-.518.541a.417.417 0 0 0 .008.588L4.913 8.4a.556.556 0 0 0 .784-.033l5.428-7.16a.394.394 0 0 0-.054-.554z" />
                <path d="M14.917.653l-.542-.45a.394.394 0 0 0-.553.063L8.961 6.505l-.498-.469a.417.417 0 0 0-.59.008l-.518.541a.417.417 0 0 0 .008.588l1.368 1.286a.556.556 0 0 0 .784-.033l5.456-7.219a.394.394 0 0 0-.054-.554z" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-[#111b21] px-2 pt-1.5 pb-2 flex items-end gap-2 shrink-0 z-40">
        <div className="bg-[#202c33] rounded-[26px] flex-1 min-h-[42px] flex items-center px-3 gap-2">
          <Smile className="w-[20px] h-[20px] text-[#8696a0] shrink-0 cursor-pointer" strokeWidth={1.8} />
          <span className="flex-1 text-[#8696a0] text-[14px] select-none cursor-text">Type a message</span>
          <Paperclip className="w-[19px] h-[19px] text-[#8696a0] shrink-0 cursor-pointer rotate-45" strokeWidth={1.8} />
          <Camera className="w-[19px] h-[19px] text-[#8696a0] shrink-0 cursor-pointer" strokeWidth={1.8} />
        </div>
        <div className="w-[42px] h-[42px] rounded-full bg-[#00a884] flex items-center justify-center shrink-0 cursor-pointer shadow-lg">
          <Mic className="w-[20px] h-[20px] text-white" strokeWidth={2} />
        </div>
      </div>

      <div className="h-5 bg-[#111b21] flex items-center justify-center shrink-0 z-40">
        <div className="w-20 h-[3px] bg-white/20 rounded-full" />
      </div>
    </>
  );
};
