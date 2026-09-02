import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

export const SplashScreen = ({ onComplete }: { onComplete: () => void }) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(onComplete, 600);
    }, 2600);
    return () => clearTimeout(timer);
  }, [onComplete]);

  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          initial={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.6, ease: 'easeInOut' }}
          className="fixed inset-0 z-[9999] bg-slate-950 flex flex-col items-center justify-center gap-8"
        >
          {/* Radial glow behind logo */}
          <motion.div
            initial={{ opacity: 0, scale: 0.5 }}
            animate={{ opacity: [0, 0.35, 0.18], scale: [0.5, 1.4, 1.1] }}
            transition={{ duration: 1.2, ease: 'easeOut', delay: 0.3 }}
            className="absolute w-56 h-56 rounded-full pointer-events-none"
            style={{
              background: 'radial-gradient(circle, rgba(37,99,235,0.65) 0%, transparent 70%)',
            }}
          />

          {/* Logo — mix-blend-mode:screen removes the white PNG background on dark bg */}
          <motion.div
            initial={{ opacity: 0, rotate: 28, scale: 0.65 }}
            animate={{ opacity: 1, rotate: 0, scale: 1 }}
            transition={{
              opacity: { duration: 0.35, delay: 0.1 },
              rotate: {
                type: 'spring',
                stiffness: 170,
                damping: 13,
                delay: 0.2,
              },
              scale: {
                type: 'spring',
                stiffness: 190,
                damping: 15,
                delay: 0.2,
              },
            }}
            className="relative w-40 h-40 z-10"
            style={{ mixBlendMode: 'screen' }}
          >
            <img
              src="/logo.png"
              alt="FlowBack"
              className="w-full h-full object-contain"
            />
          </motion.div>

          {/* Wordmark */}
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 1.1 }}
            className="flex flex-col items-center gap-1.5 z-10"
          >
            <span className="text-white text-2xl font-bold tracking-wide">
              FlowBack
            </span>
            <span className="text-slate-500 text-xs tracking-[0.25em] uppercase font-medium">
              Intelligent Payment Recovery
            </span>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};
