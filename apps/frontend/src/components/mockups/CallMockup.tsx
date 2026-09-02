import { useState, useRef, useEffect } from 'react';
import { Phone, PhoneOff, MicOff, Volume2, Plus, Grip } from 'lucide-react';
import { motion } from 'framer-motion';

export const CallMockup = ({ customer, pendingAction }: any) => {
  const [callState, setCallState] = useState<'incoming' | 'active'>('incoming');
  const [timer, setTimer] = useState(0);
  const [blobUrl, setBlobUrl] = useState<string | undefined>(undefined);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  let payload = pendingAction?.draft_payload?.RawMessage || {};
  if (!payload.body && pendingAction?.content?.String) {
    try {
      payload = JSON.parse(pendingAction.content.String);
    } catch(e) {}
  }

  const audioUrl = payload.audio_url;
  const scriptBody = payload.body;

  useEffect(() => {
    if (!audioUrl) return;

    if (audioUrl.startsWith('http')) {
      // If it's already a hosted URL, just use it directly!
      setBlobUrl(audioUrl);
      return;
    }

    try {
      let b64Data = audioUrl;
      let contentType = 'audio/wav';
      
      if (audioUrl.startsWith('data:')) {
         const parts = audioUrl.split(',');
         const match = parts[0].match(/:(.*?);/);
         if (match) contentType = match[1];
         b64Data = parts[1];
      }

      const byteCharacters = atob(b64Data);
      const byteNumbers = new Array(byteCharacters.length);
      for (let i = 0; i < byteCharacters.length; i++) {
        byteNumbers[i] = byteCharacters.charCodeAt(i);
      }
      const byteArray = new Uint8Array(byteNumbers);
      const blob = new Blob([byteArray], { type: contentType });
      const url = URL.createObjectURL(blob);
      setBlobUrl(url);
      
      return () => URL.revokeObjectURL(url);
    } catch (e) {
      let finalSrc = audioUrl;
      if (!finalSrc.startsWith('http') && !finalSrc.startsWith('data:')) {
        finalSrc = `data:audio/wav;base64,${finalSrc}`;
      }
      setBlobUrl(finalSrc);
    }
  }, [audioUrl]);

  useEffect(() => {
    let interval: any;
    if (callState === 'active') {
      interval = setInterval(() => {
        setTimer((t) => t + 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [callState]);

  const handleAccept = () => {
    if (audioRef.current && blobUrl) {
      audioRef.current.play().catch(e => {
        console.error("Audio play error:", e);
        
        // Fallback to a test beep only if strictly necessary
        if (e.name === 'NotSupportedError' || e.message.includes('not suitable')) {
           console.log("Playing fallback test audio because the backend provided an invalid/dummy audio file.");
           if (audioRef.current) {
              audioRef.current.src = "https://assets.mixkit.co/active_storage/sfx/2869/2869-preview.mp3";
              audioRef.current.play().catch(err => console.error("Fallback also failed", err));
           }
        }
      });
    }
    setCallState('active');
  };

  const handleEnd = () => {
    setCallState('incoming');
    setTimer(0);
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current.currentTime = 0;
    }
  };

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s < 10 ? '0' : ''}${s}`;
  };

  return (
    <>
      <div className="h-8 w-full bg-transparent flex justify-between items-end px-4 pb-1 text-[10.5px] font-bold text-white shrink-0 z-50 absolute top-0 left-0 right-0">
        <span>9:41</span>
        <div className="flex items-center gap-1">
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><rect x="1" y="14" width="3" height="7" rx="0.5"/><rect x="6" y="10" width="3" height="11" rx="0.5"/><rect x="11" y="6" width="3" height="15" rx="0.5"/><rect x="16" y="2" width="3" height="19" rx="0.5" opacity="0.3"/></svg>
          <svg viewBox="0 0 24 24" width="12" height="12" fill="white"><path d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2C16.93 2.93 7.08 2.93 1 9zm8 8l3 3 3-3c-1.65-1.66-4.34-1.66-6 0zm-4-4l2 2c2.76-2.76 7.24-2.76 10 0l2-2C15.14 9.14 8.87 9.14 5 13z"/></svg>
          <svg viewBox="0 0 24 24" width="15" height="15" fill="none"><rect x="1.5" y="7" width="17" height="10" rx="2" stroke="white" strokeWidth="1.5"/><rect x="3" y="8.5" width="13" height="7" rx="1" fill="white"/><path d="M20 10.5v3" stroke="white" strokeWidth="1.5" strokeLinecap="round"/></svg>
        </div>
      </div>

      <div className="absolute inset-0 flex flex-col items-center pt-24 z-40 rounded-[2rem] overflow-hidden bg-gradient-to-b from-[#1a1c29] to-[#0a0a0f]">
        
        {/* Caller Info */}
        <div className="flex flex-col items-center w-full px-6">
          <div className="w-20 h-20 rounded-full bg-slate-700 flex items-center justify-center text-3xl font-light text-white mb-4 shadow-lg">
             {customer?.name?.String?.slice(0, 1).toUpperCase() || 'C'}
          </div>
          <div className="text-white text-2xl font-light tracking-wide text-center">
             {customer?.name?.String || 'Customer'}
          </div>
          <div className="text-white/60 text-sm mt-2 font-medium">
            {callState === 'incoming' ? 'FlowBack Support Call' : formatTime(timer)}
          </div>
        </div>

        {callState === 'active' && scriptBody && (
           <div className="mt-8 px-6 w-full flex-1 overflow-hidden flex flex-col">
              <div className="bg-white/5 border border-white/10 rounded-xl p-4 text-white/80 text-xs leading-relaxed italic overflow-y-auto mb-32">
                 "{scriptBody.replace(/\[.*?\]\s*/g, '')}"
              </div>
           </div>
        )}

        {/* Action Buttons */}
        <div className="absolute bottom-12 left-0 right-0 px-8 flex flex-col items-center">
           {callState === 'incoming' ? (
             <div className="flex justify-between w-full px-6">
               <div className="flex flex-col items-center gap-2">
                 <button className="w-16 h-16 rounded-full bg-[#ff3b30] flex items-center justify-center shadow-lg transition-transform hover:scale-105 active:scale-95">
                   <PhoneOff className="w-8 h-8 text-white" />
                 </button>
                 <span className="text-white/80 text-xs">Decline</span>
               </div>
               
               <div className="flex flex-col items-center gap-2">
                 <motion.button 
                   onClick={handleAccept}
                   animate={{ scale: [1, 1.1, 1], boxShadow: ["0 0 0px rgba(52,199,89,0)", "0 0 20px rgba(52,199,89,0.6)", "0 0 0px rgba(52,199,89,0)"] }}
                   transition={{ repeat: Infinity, duration: 2 }}
                   className="w-16 h-16 rounded-full bg-[#34c759] flex items-center justify-center shadow-lg transition-transform hover:scale-105 active:scale-95"
                 >
                   <Phone className="w-8 h-8 text-white fill-current" />
                 </motion.button>
                 <span className="text-white/80 text-xs">Accept</span>
               </div>
             </div>
           ) : (
             <div className="flex flex-col w-full items-center gap-8 bg-[#0a0a0f] pt-4">
                {/* In call controls */}
                <div className="grid grid-cols-3 gap-y-6 gap-x-8 w-full max-w-[240px]">
                  {[
                    { icon: MicOff, label: 'mute' },
                    { icon: Grip, label: 'keypad' },
                    { icon: Volume2, label: 'speaker' },
                    { icon: Plus, label: 'add call' },
                    { icon: Phone, label: 'FaceTime' },
                    { icon: Phone, label: 'contacts' },
                  ].map((btn, i) => (
                    <div key={i} className="flex flex-col items-center gap-2 opacity-60">
                      <div className="w-14 h-14 rounded-full bg-white/20 flex items-center justify-center">
                        <btn.icon className="w-6 h-6 text-white" />
                      </div>
                      <span className="text-white/80 text-[10px]">{btn.label}</span>
                    </div>
                  ))}
                </div>
                
                <button onClick={handleEnd} className="w-16 h-16 rounded-full bg-[#ff3b30] flex items-center justify-center shadow-lg transition-transform hover:scale-105 active:scale-95">
                  <PhoneOff className="w-8 h-8 text-white" />
                </button>
             </div>
           )}
        </div>

        {blobUrl && (
          <audio ref={audioRef} src={blobUrl} className="hidden" preload="auto" />
        )}
      </div>
    </>
  );
};
