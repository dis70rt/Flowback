import { useEffect } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { Home, List, Users } from 'lucide-react';
import { useMetrics } from '../../hooks/useApi';
import { Skeleton } from '@/components/ui/skeleton';
import { Toaster } from '@/components/ui/sonner';
import { toast } from 'sonner';
import { useQueryClient } from '@tanstack/react-query';

export const WorkspaceLayout = () => {
  const { data: metrics, isPending } = useMetrics();
  const queryClient = useQueryClient();

  useEffect(() => {
    const sse = new EventSource('/api/stream');
    
    sse.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        toast('New Activity', {
          description: `Case ${data.case_id?.slice(0, 8)}: ${data.event}`,
        });
        // Refresh queues and metrics
        queryClient.invalidateQueries({ queryKey: ['cases'] });
        queryClient.invalidateQueries({ queryKey: ['metrics'] });
      } catch (e) {
        console.error('SSE parse error:', e);
      }
    };

    return () => sse.close();
  }, [queryClient]);

  return (
    <div className="dark flex h-screen bg-slate-950 text-white overflow-hidden">
      <Toaster theme="dark" />
      {/* Sidebar */}
      <aside className="w-64 border-r border-slate-800 p-4 flex flex-col gap-6">
        <div className="flex items-center gap-2 px-2">
          <div className="w-8 h-8 bg-purple-900/50 text-purple-300 border border-purple-500/30 rounded flex items-center justify-center font-bold">F</div>
          <h1 className="font-semibold text-lg">FlowBack</h1>
        </div>

        <nav className="flex flex-col gap-2">
          <NavLink to="/" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
            <Home className="w-4 h-4" /> Overview
          </NavLink>
          <NavLink to="/cases" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
            <List className="w-4 h-4" /> Live Queue
          </NavLink>
          <NavLink to="/customers" className={({ isActive }) => `flex items-center gap-3 px-3 py-2 rounded-md transition-colors ${isActive ? 'bg-slate-800 text-white' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}`}>
            <Users className="w-4 h-4" /> Customers
          </NavLink>
        </nav>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col h-full overflow-hidden">
        {/* Top Header Metrics */}
        <header className="h-20 border-b border-slate-800 flex items-center justify-between px-8">
          <div className="flex items-center gap-8">
            {isPending ? (
              <div className="flex gap-4">
                 <Skeleton className="h-12 w-48 bg-slate-800" />
                 <Skeleton className="h-12 w-32 bg-slate-800" />
                 <Skeleton className="h-12 w-32 bg-slate-800" />
              </div>
            ) : (
              <>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-purple-900/30 text-purple-300 border border-purple-500/20 flex items-center justify-center">₹</div>
                  <div>
                    <div className="text-xs text-slate-400">Total Revenue Recovered</div>
                    <div className="font-semibold text-xl">${((metrics?.total_revenue_recovered || 0) / 100).toLocaleString()}</div>
                  </div>
                </div>
                
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-cyan-900/30 text-cyan-300 border border-cyan-500/20 flex items-center justify-center"><List className="w-4 h-4" /></div>
                  <div>
                    <div className="text-xs text-slate-400">Active Cases</div>
                    <div className="font-semibold text-xl">{metrics?.active_cases || 0}</div>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-emerald-900/30 text-emerald-300 border border-emerald-500/20 flex items-center justify-center">✓</div>
                  <div>
                    <div className="text-xs text-slate-400">AI Success Rate</div>
                    <div className="font-semibold text-xl">{metrics?.ai_success_rate || 0}%</div>
                  </div>
                </div>
              </>
            )}
          </div>
          
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-slate-800 border border-slate-700 text-slate-300 flex items-center justify-center text-xs font-bold">AD</div>
            <div className="text-sm">
              <div className="font-medium">Admin</div>
              <div className="text-xs text-slate-400">Supervisor</div>
            </div>
          </div>
        </header>

        {/* Dynamic Screen Content */}
        <div className="flex-1 overflow-hidden">
          <Outlet />
        </div>
      </main>
    </div>
  );
};
