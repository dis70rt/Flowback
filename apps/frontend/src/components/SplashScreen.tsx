import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

export const SplashScreen = ({ onComplete }: { onComplete: () => void }) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(onComplete, 700);
    }, 3000);
    return () => clearTimeout(timer);
  }, [onComplete]);

  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          initial={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.7, ease: 'easeInOut' }}
          className="fixed inset-0 z-[9999] bg-slate-950 flex flex-col items-center justify-center gap-10"
        >
          {/* Blue glow blooms after logo is revealed */}
          <motion.div
            className="absolute rounded-full pointer-events-none"
            style={{
              width: 320,
              height: 320,
              background:
                'radial-gradient(circle, rgba(59,130,246,0.45) 0%, transparent 70%)',
            }}
            initial={{ opacity: 0, scale: 0.7 }}
            animate={{ opacity: [0, 0, 0.8, 0.3], scale: [0.7, 0.7, 1.2, 1.0] }}
            transition={{
              duration: 2.8,
              times: [0, 0.45, 0.7, 1.0],
              ease: 'easeOut',
            }}
          />

          {/* Logo — diagonal wipe from bottom-left (arrow) → top-right (S-tip)
              mix-blend-mode:screen removes the white PNG background on slate-950 */}
          <div className="relative" style={{ width: 220, height: 180 }}>
            <motion.div
              className="absolute inset-0"
              style={{ mixBlendMode: 'screen' }}
              initial={{
                // A degenerate polygon at the bottom-left: nothing visible
                clipPath: 'polygon(0% 0%, 0% 100%, 0% 100%, 0% 0%)',
              }}
              animate={{
                // The right edge of the revealed area sweeps diagonally from BL to TR.
                // The leading diagonal edge: bottom goes to 190%, top goes to 90%,
                // so the sweep angle matches the logo's own diagonal (bottom-left arrow → top-right S-tip).
                clipPath: 'polygon(0% 0%, 0% 100%, 190% 100%, 90% 0%)',
              }}
              transition={{
                duration: 1.3,
                ease: [0.25, 0.46, 0.45, 0.94],
                delay: 0.15,
              }}
            >
              <img
                src="/logo.png"
                alt="FlowBack"
                className="w-full h-full object-contain"
                draggable={false}
              />
            </motion.div>
          </div>

          {/* Wordmark fades in after the logo is fully revealed */}
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 1.65 }}
            className="flex flex-col items-center gap-2 relative"
          >
            <span className="text-white text-[1.75rem] font-bold tracking-wide leading-none">
              FlowBack
            </span>
            <span className="text-slate-500 text-[0.65rem] tracking-[0.28em] uppercase font-semibold">
              Intelligent Payment Recovery
            </span>
          </motion.div>

          {/* Thin loading bar — fills over the full reveal duration */}
          <div className="absolute bottom-10 w-36 h-[2px] bg-slate-800 rounded-full overflow-hidden">
            <motion.div
              className="h-full rounded-full"
              style={{ background: 'linear-gradient(90deg, #1d4ed8, #3b82f6, #93c5fd)' }}
              initial={{ width: '0%' }}
              animate={{ width: '100%' }}
              transition={{ duration: 2.3, delay: 0.15, ease: 'easeInOut' }}
            />
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
