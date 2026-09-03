import { useState, useEffect } from 'react';
import { useCases, useCaseDetails, useCustomer, useCustomers, useCustomerPayments, useCustomerCommunications } from '../hooks/useApi';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Search, User, Mail, Phone, MapPin, CreditCard, MessageSquare, Clock, ShieldCheck, } from 'lucide-react';

export const Customers = () => {
  const [searchInput, setSearchInput] = useState('');
  const [activeCustomerId, setActiveCustomerId] = useState<string | null>(null);

  const { data: customersList, isPending: isLoadingCustomersList } = useCustomers();

  // Auto-select a customer for demo purposes if none is selected
  const { data: casesData } = useCases(1, 1);
  const firstCaseId = casesData?.data?.[0]?.id;
  const { data: firstCaseDetails } = useCaseDetails(firstCaseId);

  useEffect(() => {
    if (!activeCustomerId && customersList && customersList.length > 0) {
      setActiveCustomerId(customersList[0].id);
    } else if (!activeCustomerId && firstCaseDetails?.case?.customer_id) {
      setActiveCustomerId(firstCaseDetails.case.customer_id);
    }
  }, [firstCaseDetails, activeCustomerId, customersList]);

  const { data: customer, isPending: isLoadingCustomer } = useCustomer(activeCustomerId || undefined);
  const { data: payments, isPending: isLoadingPayments } = useCustomerPayments(activeCustomerId || undefined);
  const { data: communications, isPending: isLoadingComms } = useCustomerCommunications(activeCustomerId || undefined);

  
  const filteredCustomers = customersList?.filter(c => 
    c.name?.String?.toLowerCase().includes(searchInput.toLowerCase()) || 
    c.email?.String?.toLowerCase().includes(searchInput.toLowerCase())
  ) || [];

  return (
    <div className="flex h-full w-full bg-slate-950 flex-col overflow-hidden">
      {/* Header & Search */}
      <div className="px-8 pt-8 pb-6 border-b border-slate-800/60 flex items-center justify-between shrink-0 bg-slate-950/80 backdrop-blur-md z-10">
        <div>
          <h1 className="text-2xl font-bold text-slate-100 tracking-tight">Customer Directory</h1>
          <p className="text-slate-500 text-sm mt-1">Search and manage customer profiles, payment history, and communications.</p>
        </div>
        <div className="relative w-96">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
          <input 
            type="text" 
            placeholder="Search by name or email..."
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            className="w-full bg-slate-900/50 border border-slate-800/80 rounded-lg pl-10 pr-4 py-2 text-sm text-slate-200 placeholder:text-slate-600 focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all"
          />
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Customer List Sidebar */}
        <div className="w-80 border-r border-slate-800/60 bg-slate-900/20 flex flex-col shrink-0">
          <div className="p-4 border-b border-slate-800/40 text-xs font-semibold text-slate-500 uppercase tracking-wider sticky top-0 bg-slate-900/90 backdrop-blur z-10">
            {filteredCustomers.length} Customers
          </div>
          <div className="flex-1 overflow-y-auto p-2 space-y-1">
            {isLoadingCustomersList ? (
              Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="p-3 rounded-lg flex gap-3">
                  <Skeleton className="w-10 h-10 rounded-full bg-slate-800" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-24 bg-slate-800" />
                    <Skeleton className="h-3 w-32 bg-slate-800/50" />
                  </div>
                </div>
              ))
            ) : filteredCustomers.map(c => (
              <button
                key={c.id}
                onClick={() => setActiveCustomerId(c.id)}
                className={`w-full text-left p-3 rounded-lg flex items-center gap-3 transition-colors ${
                  activeCustomerId === c.id 
                    ? 'bg-indigo-500/10 border border-indigo-500/20' 
                    : 'hover:bg-slate-800/40 border border-transparent'
                }`}
              >
                <div className={`w-10 h-10 rounded-full flex items-center justify-center shrink-0 ${
                  activeCustomerId === c.id ? 'bg-indigo-500/20 text-indigo-400' : 'bg-slate-800 text-slate-400'
                }`}>
                  <User className="w-5 h-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className={`font-medium text-sm truncate ${activeCustomerId === c.id ? 'text-indigo-200' : 'text-slate-300'}`}>
                    {c.name?.String || 'Unknown Customer'}
                  </div>
                  <div className="text-xs text-slate-500 truncate mt-0.5">
                    {c.email?.String || c.phone?.String || 'No contact info'}
                  </div>
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Customer Detail Content */}
        <div className="flex-1 overflow-y-auto p-8">
        {!activeCustomerId ? (
          <div className="h-full flex flex-col items-center justify-center text-slate-500">
            <User className="w-12 h-12 mb-4 text-slate-700" />
            <h2 className="text-lg font-medium text-slate-300">No Customer Selected</h2>
            <p className="text-sm mt-1 text-slate-500 max-w-md text-center">
              Select a customer from the sidebar to view their complete profile, payment history, and communication logs.
            </p>
          </div>
        ) : isLoadingCustomer ? (
          <div className="space-y-6">
            <Skeleton className="h-48 w-full bg-slate-800/50 rounded-2xl" />
            <div className="grid grid-cols-2 gap-6">
              <Skeleton className="h-64 w-full bg-slate-800/50 rounded-2xl" />
              <Skeleton className="h-64 w-full bg-slate-800/50 rounded-2xl" />
            </div>
          </div>
        ) : (
          <div className="max-w-6xl mx-auto space-y-6 pb-12">
            
            {/* Top Profile Card */}
            <div className="bg-slate-900/40 border border-slate-800/80 rounded-2xl p-8 backdrop-blur-sm relative overflow-hidden">
              <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500/5 rounded-full blur-3xl -mr-20 -mt-20 pointer-events-none" />
              
              <div className="flex items-start justify-between relative z-10">
                <div className="flex items-center gap-6">
                  <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-indigo-900/80 to-slate-800 text-indigo-300 flex items-center justify-center text-3xl font-bold border border-indigo-500/20 shadow-xl shadow-indigo-900/20">
                    {customer?.name?.String?.slice(0, 2).toUpperCase() || 'C'}
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold text-slate-100 tracking-tight">{customer?.name?.String || 'Unknown Customer'}</h2>
                    <div className="flex items-center gap-3 mt-2 text-sm text-slate-400">
                      <span className="font-mono bg-slate-950/50 px-2 py-0.5 rounded border border-slate-800">ID: {customer?.id}</span>
                      {customer?.value_tier?.Valid && (
                        <Badge variant="outline" className="text-indigo-300 border-indigo-500/30 bg-indigo-500/10">
                          {customer.value_tier.String} Tier
                        </Badge>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex gap-8">
                  <div className="text-right">
                    <div className="text-xs uppercase tracking-wider text-slate-500 font-semibold mb-1">Reliability Score</div>
                    <div className="flex items-baseline justify-end gap-1">
                      <span className="text-3xl font-bold text-emerald-400">{customer?.reliability_score || 0}</span>
                      <span className="text-slate-600 text-sm">/100</span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs uppercase tracking-wider text-slate-500 font-semibold mb-1">Payments</div>
                    <div className="flex items-baseline justify-end gap-1">
                      <span className="text-3xl font-bold text-slate-200">{customer?.successful_payments || 0}</span>
                      <span className="text-slate-600 text-sm">/ {customer?.total_payments || 0}</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-4 gap-6 mt-8 pt-6 border-t border-slate-800/60 relative z-10">
                <div className="flex items-center gap-3 text-sm text-slate-300">
                  <div className="w-8 h-8 rounded-full bg-slate-800/80 flex items-center justify-center shrink-0">
                    <Mail className="w-4 h-4 text-slate-400" />
                  </div>
                  <span className="truncate">{customer?.email?.String || 'No email provided'}</span>
                </div>
                <div className="flex items-center gap-3 text-sm text-slate-300">
                  <div className="w-8 h-8 rounded-full bg-slate-800/80 flex items-center justify-center shrink-0">
                    <Phone className="w-4 h-4 text-slate-400" />
                  </div>
                  <span className="truncate">{customer?.phone?.String || 'No phone provided'}</span>
                </div>
                <div className="flex items-center gap-3 text-sm text-slate-300">
                  <div className="w-8 h-8 rounded-full bg-slate-800/80 flex items-center justify-center shrink-0">
                    <MapPin className="w-4 h-4 text-slate-400" />
                  </div>
                  <span className="truncate">{[customer?.city?.String, customer?.state?.String].filter(Boolean).join(', ') || 'No location'}</span>
                </div>
                <div className="flex items-center gap-3 text-sm text-slate-300">
                  <div className="w-8 h-8 rounded-full bg-slate-800/80 flex items-center justify-center shrink-0">
                    <MessageSquare className="w-4 h-4 text-slate-400" />
                  </div>
                  <span className="truncate capitalize text-indigo-200">{customer?.preferred_channel?.String?.replace('send_', '') || 'Unknown channel'}</span>
                </div>
              </div>
            </div>

            {/* Bottom Two Columns */}
            <div className="grid grid-cols-2 gap-6">
              
              {/* Payment History */}
              <div className="bg-slate-900/40 border border-slate-800/80 rounded-2xl p-6 backdrop-blur-sm flex flex-col min-h-[400px]">
                <div className="flex items-center gap-2 mb-6">
                  <CreditCard className="w-5 h-5 text-indigo-400" />
                  <h3 className="text-lg font-medium text-slate-200">Payment History</h3>
                </div>
                
                {isLoadingPayments ? (
                  <div className="space-y-4">
                    {[1,2,3].map(i => <Skeleton key={i} className="h-16 w-full bg-slate-800/50 rounded-lg" />)}
                  </div>
                ) : (
                  <div className="flex-1 flex flex-col items-center justify-center text-center p-6 border border-dashed border-slate-700/50 rounded-xl bg-slate-950/30">
                    <ShieldCheck className="w-10 h-10 text-slate-600 mb-3" />
                    <p className="text-slate-400 font-medium">{payments?.message || 'Payment history will appear here'}</p>
                    <p className="text-sm text-slate-500 mt-1">Integration with Razorpay is pending setup.</p>
                  </div>
                )}
              </div>

              {/* Communication Logs */}
              <div className="bg-slate-900/40 border border-slate-800/80 rounded-2xl p-6 backdrop-blur-sm flex flex-col min-h-[400px]">
                <div className="flex items-center gap-2 mb-6">
                  <MessageSquare className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-medium text-slate-200">Communication Logs</h3>
                </div>
                
                {isLoadingComms ? (
                  <div className="space-y-4">
                    {[1,2,3].map(i => <Skeleton key={i} className="h-16 w-full bg-slate-800/50 rounded-lg" />)}
                  </div>
                ) : communications?.length === 0 ? (
                  <div className="flex-1 flex flex-col items-center justify-center text-center p-6 border border-dashed border-slate-700/50 rounded-xl bg-slate-950/30">
                    <MessageSquare className="w-10 h-10 text-slate-600 mb-3" />
                    <p className="text-slate-400 font-medium">No previous communications</p>
                    <p className="text-sm text-slate-500 mt-1">This customer has not received any recovery messages yet.</p>
                  </div>
                ) : (
                  <div className="space-y-3 overflow-y-auto pr-2 -mr-2">
                    {communications?.map((comm: any) => (
                      <div key={comm.id} className="bg-slate-950/60 border border-slate-800/60 rounded-xl p-4 flex gap-4 items-start">
                        <div className="w-8 h-8 rounded-full bg-slate-900 flex items-center justify-center shrink-0 mt-0.5 border border-slate-800">
                          {comm.channel === 'whatsapp' ? (
                            <MessageSquare className="w-4 h-4 text-emerald-400" />
                          ) : comm.channel === 'sms' ? (
                            <MessageSquare className="w-4 h-4 text-indigo-400" />
                          ) : (
                            <Mail className="w-4 h-4 text-slate-400" />
                          )}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center justify-between mb-1">
                            <span className="font-medium text-slate-200 text-sm capitalize">{comm.channel} Message</span>
                            <span className="text-[11px] text-slate-500 flex items-center gap-1">
                              <Clock className="w-3 h-3" /> 
                              {new Date(comm.sent_at).toLocaleDateString()}
                            </span>
                          </div>
                          <div className="flex items-center gap-2">
                            <Badge variant="outline" className={`text-[10px] px-1.5 py-0 ${
                              comm.status === 'delivered' ? 'text-emerald-400 border-emerald-400/20 bg-emerald-400/10' :
                              comm.status === 'failed' ? 'text-rose-400 border-rose-400/20 bg-rose-400/10' :
                              'text-amber-400 border-amber-400/20 bg-amber-400/10'
                            }`}>
                              {comm.status}
                            </Badge>
                            <span className="text-[11px] text-slate-500 font-mono truncate">SID: {comm.message_sid?.String || 'N/A'}</span>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

            </div>
          </div>
        )}
        </div>
      </div>
    </div>
  );
};
