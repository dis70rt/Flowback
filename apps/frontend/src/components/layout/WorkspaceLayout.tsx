import { NavLink, Outlet } from 'react-router-dom';
import { Home, List, Users, HelpCircle } from 'lucide-react';
import { Toaster } from '@/components/ui/sonner';
import { LiveActionsProvider } from '../../hooks/useLiveActions';

export const WorkspaceLayout = () => {
  return (
    <LiveActionsProvider>
      <div className="dark flex h-screen bg-slate-950 text-white overflow-hidden">
        <Toaster theme="dark" />
        
        {/* Sidebar */}
        <aside className="w-64 border-r border-slate-800 flex flex-col shrink-0">
          <div className="p-4 flex flex-col gap-6 flex-1">
            <div className="flex items-center gap-2 px-2 mt-2">
              <img src="/logo.png" alt="FlowBack Logo" className="w-8 h-8 object-contain" />
              <h1 className="font-semibold text-lg tracking-tight">FlowBack</h1>
            </div>

            <nav className="flex flex-col gap-1.5 mt-4">
              <NavLink to="/" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
                <Home className="w-4 h-4" /> Overview
              </NavLink>
              <NavLink to="/cases" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
                <List className="w-4 h-4" /> Live Queue
              </NavLink>
              <NavLink to="/customers" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
                <Users className="w-4 h-4" /> Customers
              </NavLink>
            </nav>
          </div>

          <div className="p-4 border-t border-slate-800/60 flex flex-col gap-2">
            <div className="flex items-center gap-3 px-3 py-2 hover:bg-slate-800/40 rounded-lg cursor-pointer transition-colors">
              <div className="w-8 h-8 rounded-full bg-slate-800 border border-slate-700 text-slate-300 flex items-center justify-center text-xs font-bold shrink-0 shadow-sm">AD</div>
              <div className="text-sm min-w-0">
                <div className="font-medium truncate text-slate-200">Admin User</div>
                <div className="text-[11px] text-slate-500 truncate">Workspace Owner</div>
              </div>
            </div>
            
            <button className="flex items-center gap-3 px-3 py-2 mt-2 rounded-md transition-colors text-sm font-medium text-slate-400 hover:text-white hover:bg-slate-800/50 w-full text-left">
              <HelpCircle className="w-4 h-4" /> Help & Support
            </button>
          </div>
        </aside>

        {/* Main Content Area */}
        <main className="flex-1 flex flex-col h-full overflow-hidden">
          {/* Dynamic Screen Content */}
          <div className="flex-1 overflow-hidden">
            <Outlet />
          </div>
        </main>
      </div>
    </LiveActionsProvider>
  );
};
