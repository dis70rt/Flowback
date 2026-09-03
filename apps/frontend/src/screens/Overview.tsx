import { useMetricsOverview, useMetricsTrends, useMetricsChannels, useMetricsPipeline, useMetricsRecovered } from '../hooks/useApi';
import { ResponsiveLine } from '@nivo/line';
import { ResponsivePie } from '@nivo/pie';
import { Skeleton } from '@/components/ui/skeleton';
import { Activity, AlertTriangle, TrendingUp, CheckCircle, ShieldAlert, ChevronDown } from 'lucide-react';

const colors = {
  bg: '#171B2A',
  cardBg: '#1D2335',
  cardBorder: '#2B344A',
  textMain: '#F1F4FC',
  textImportant: '#F1F4FC',
  textSecondary: '#AAB4CC',
  textLabel: '#7F8DAA',
  textAxisTable: '#8290AC',
  divider: '#293147',
  hover: '#222A3D',
  kpi: {
    atRisk: '#FF647C',
    recovered: '#6FCF5B',
    active: '#4F8CFF',
    ai: '#A78BFA'
  },
  channels: {
    email: '#5B8DEF',
    sms: '#A78BFA',
    call: '#F3B95F',
    whatsapp: '#5CC8FF'
  }
};

const nivoTheme = {
  background: 'transparent',
  textColor: colors.textAxisTable,
  fontSize: 11,
  axis: {
    domain: {
      line: { stroke: 'transparent', strokeWidth: 1 }
    },
    legend: {
      text: { fill: colors.textAxisTable }
    },
    ticks: {
      line: { stroke: 'transparent', strokeWidth: 1 },
      text: { fill: colors.textAxisTable }
    }
  },
  grid: {
    line: { stroke: colors.divider, strokeWidth: 1 }
  },
  legends: {
    text: { fill: colors.textSecondary }
  },
  tooltip: {
    container: { background: colors.cardBg, color: colors.textMain, border: `1px solid ${colors.cardBorder}` }
  }
};

const formatCurrencyAbbr = (valPaise: number) => {
  const val = valPaise / 100;
  if (val >= 10000000) return `₹${(val / 10000000).toFixed(1)}Cr`;
  if (val >= 100000) return `₹${(val / 100000).toFixed(1)}L`;
  if (val >= 1000) return `₹${(val / 1000).toFixed(1)}K`;
  return `₹${val}`;
};

// Custom layer for center metric in pie chart
const CenteredMetric = ({ dataWithArc, centerX, centerY }: any) => {
  let total = 0;
  dataWithArc.forEach((datum: any) => {
    total += datum.value;
  });
  return (
    <g transform={`translate(${centerX}, ${centerY})`}>
      <text
        textAnchor="middle"
        dominantBaseline="central"
        style={{ fontSize: '24px', fontWeight: 600, fill: colors.textMain }}
      >
        {new Intl.NumberFormat().format(total)}
      </text>
      <text
        textAnchor="middle"
        dominantBaseline="central"
        y={20}
        style={{ fontSize: '12px', fill: colors.textSecondary }}
      >
        Interactions
      </text>
    </g>
  );
};

export const Overview = () => {
  const { data: metrics, isPending: loadingMetrics } = useMetricsOverview();
  const { data: trends, isPending: loadingTrends } = useMetricsTrends();
  const { data: channels, isPending: loadingChannels } = useMetricsChannels();
  const { data: recovered, isPending: loadingRecovered } = useMetricsRecovered();

  const trendData = [
    {
      id: 'Recovered',
      color: colors.kpi.recovered,
      data: (trends || []).map(t => ({ x: new Date(t.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }), y: t.daily_recovered / 100 }))
    },
    {
      id: 'Failed',
      color: colors.kpi.atRisk,
      data: (trends || []).map(t => ({ x: new Date(t.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }), y: t.daily_failed / 100 }))
    }
  ];

  // Calculate percentages for channel legend
  const totalChannels = (channels || []).reduce((acc, c) => acc + c.count, 0);
  
  const channelData = (channels || []).map(c => {
    let label = c.channel.String.replace('send_', '');
    label = label.charAt(0).toUpperCase() + label.slice(1);
    
    // Assign specific colors requested
    let color = colors.channels.email;
    if (label.toLowerCase() === 'sms') color = colors.channels.sms;
    if (label.toLowerCase() === 'call') color = colors.channels.call;
    if (label.toLowerCase() === 'whatsapp') color = colors.channels.whatsapp;
    
    const pct = totalChannels > 0 ? Math.round((c.count / totalChannels) * 100) : 0;
    
    return {
      id: label,
      label: `${label} ${pct}%`,
      value: c.count,
      color: color
    };
  });

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto" style={{ backgroundColor: colors.bg }}>
      <div className="px-8 pt-8 pb-6 border-b shrink-0 z-10" style={{ borderColor: colors.cardBorder, backgroundColor: colors.bg }}>
        <div>
          <h1 className="text-[28px] font-bold tracking-tight" style={{ color: colors.textMain }}>Overview Analytics</h1>
          <p className="text-sm mt-1" style={{ color: colors.textSecondary }}>Business metrics, recovery performance, and AI intelligence.</p>
        </div>
      </div>

      <div className="flex-1 p-8 space-y-[28px] w-full">
        {/* Top Metric Cards */}
        <div className="grid grid-cols-4 gap-7">
          {[
            { label: 'AMOUNT AT RISK', value: metrics ? formatCurrencyAbbr(metrics.total_amount_at_risk || 0) : '0', delta: '↑ 12.4% vs previous 30 days', icon: AlertTriangle, iconColor: colors.kpi.atRisk },
            { label: 'AMOUNT RECOVERED', value: metrics ? formatCurrencyAbbr(metrics.total_amount_recovered || 0) : '0', delta: '↑ 8.2% vs previous 30 days', icon: CheckCircle, iconColor: colors.kpi.recovered },
            { label: 'ACTIVE CASES', value: metrics?.active_cases !== undefined ? metrics.active_cases : '0', delta: '↓ 3.1% vs previous 30 days', icon: Activity, iconColor: colors.kpi.active },
            { label: 'AI SUCCESS RATE', value: metrics?.ai_success_rate !== undefined ? `${metrics.ai_success_rate}%` : '0%', delta: '↑ 1.5% vs previous 30 days', icon: TrendingUp, iconColor: colors.kpi.ai },
          ].map((card, idx) => (
            <div key={idx} className="rounded-xl p-5 border flex flex-col justify-between h-[115px]" style={{ backgroundColor: colors.cardBg, borderColor: colors.cardBorder }}>
              <div className="flex items-center justify-between mb-2">
                <span className="text-[11px] font-semibold tracking-wider" style={{ color: colors.textLabel }}>{card.label}</span>
                <card.icon className="w-4 h-4" style={{ color: card.iconColor }} />
              </div>
              <div>
                <div className="text-[32px] font-bold leading-none mb-1.5" style={{ color: colors.textImportant }}>
                  {loadingMetrics ? <Skeleton className="h-8 w-24" style={{ backgroundColor: colors.divider }} /> : card.value}
                </div>
                <div className="text-[12px]" style={{ color: colors.textSecondary }}>
                  {card.delta}
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Charts Row */}
        <div className="grid grid-cols-[2.1fr_1fr] gap-7 h-[400px]">
          {/* Line Chart */}
          <div className="rounded-xl border p-6 flex flex-col" style={{ backgroundColor: colors.cardBg, borderColor: colors.cardBorder }}>
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-base font-semibold" style={{ color: colors.textMain }}>Recovery Trends</h3>
              <button className="flex items-center text-sm gap-1.5 px-3 py-1.5 rounded-md transition-colors" style={{ color: colors.textSecondary, backgroundColor: colors.hover }}>
                Last 30 days <ChevronDown className="w-3 h-3" />
              </button>
            </div>
            <div className="flex-1 min-h-0">
              {loadingTrends ? (
                <Skeleton className="w-full h-full rounded-lg" style={{ backgroundColor: colors.divider }} />
              ) : trendData[0].data.length > 0 ? (
                <ResponsiveLine
                  data={trendData}
                  theme={nivoTheme}
                  colors={d => d.color}
                  margin={{ top: 10, right: 20, bottom: 30, left: 55 }}
                  xScale={{ type: 'point' }}
                  yScale={{ type: 'linear', min: 0, max: 'auto', stacked: false, reverse: false }}
                  yFormat=" >-.0f"
                  curve="monotoneX"
                  lineWidth={2}
                  enableArea={true}
                  areaOpacity={0.05}
                  axisTop={null}
                  axisRight={null}
                  axisBottom={{
                    tickSize: 0,
                    tickPadding: 12,
                    tickRotation: 0,
                    tickValues: 6, // Ask Nivo to show ~6 ticks
                  }}
                  axisLeft={{
                    tickSize: 0,
                    tickPadding: 10,
                    tickRotation: 0,
                    tickValues: 5,
                    format: (value) => formatCurrencyAbbr(value * 100) // Re-multiply by 100 for formatter
                  }}
                  enablePoints={false}
                  enableGridX={false}
                  gridYValues={5}
                  useMesh={true}
                  legends={[
                    {
                      anchor: 'top-right',
                      direction: 'row',
                      justify: false,
                      translateX: 10,
                      translateY: -35,
                      itemsSpacing: 10,
                      itemDirection: 'left-to-right',
                      itemWidth: 90,
                      itemHeight: 20,
                      itemOpacity: 0.85,
                      symbolSize: 8,
                      symbolShape: 'circle',
                    }
                  ]}
                />
              ) : (
                <div className="flex items-center justify-center h-full text-sm" style={{ color: colors.textSecondary }}>No trend data available.</div>
              )}
            </div>
          </div>

          {/* Pie Chart */}
          <div className="rounded-xl border p-6 flex flex-col" style={{ backgroundColor: colors.cardBg, borderColor: colors.cardBorder }}>
            <h3 className="text-base font-semibold mb-2" style={{ color: colors.textMain }}>Channel Distribution</h3>
            <div className="flex-1 min-h-0 relative">
              {loadingChannels ? (
                <Skeleton className="w-full h-full rounded-full" style={{ backgroundColor: colors.divider }} />
              ) : channelData.length > 0 ? (
                <ResponsivePie
                  data={channelData}
                  theme={nivoTheme}
                  colors={{ datum: 'data.color' }}
                  margin={{ top: 20, right: 140, bottom: 20, left: 20 }}
                  innerRadius={0.75}
                  padAngle={1.5}
                  cornerRadius={4}
                  activeOuterRadiusOffset={4}
                  borderWidth={0}
                  enableArcLabels={false}
                  enableArcLinkLabels={false}
                  layers={['arcs', 'arcLabels', 'arcLinkLabels', 'legends', CenteredMetric]}
                  legends={[
                    {
                      anchor: 'right',
                      direction: 'column',
                      justify: false,
                      translateX: 110,
                      translateY: 0,
                      itemsSpacing: 8,
                      itemWidth: 100,
                      itemHeight: 18,
                      itemTextColor: colors.textSecondary,
                      itemDirection: 'left-to-right',
                      itemOpacity: 1,
                      symbolSize: 10,
                      symbolShape: 'circle'
                    }
                  ]}
                />
              ) : (
                <div className="flex items-center justify-center h-full text-sm" style={{ color: colors.textSecondary }}>No channel data available.</div>
              )}
            </div>
          </div>
        </div>

        {/* Data Table Row */}
        <div className="rounded-xl border flex flex-col overflow-hidden" style={{ backgroundColor: colors.cardBg, borderColor: colors.cardBorder }}>
          <div className="px-6 py-5 border-b" style={{ borderColor: colors.divider }}>
            <h3 className="text-base font-semibold" style={{ color: colors.textMain }}>Recent Recoveries</h3>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left border-collapse">
              <thead className="text-[13px] border-b" style={{ color: colors.textAxisTable, borderColor: colors.divider }}>
                <tr>
                  <th className="px-6 py-4 font-medium">Customer</th>
                  <th className="px-6 py-4 font-medium">Tier</th>
                  <th className="px-6 py-4 font-medium">Amount</th>
                  <th className="px-6 py-4 font-medium">Channel</th>
                  <th className="px-6 py-4 font-medium">Recovered Date</th>
                </tr>
              </thead>
              <tbody style={{ color: '#D8DEEB' }}>
                {loadingRecovered ? (
                  Array.from({ length: 4 }).map((_, i) => (
                    <tr key={i} className="border-b" style={{ borderColor: colors.divider }}>
                      <td className="px-6 py-4"><Skeleton className="h-4 w-32" style={{ backgroundColor: colors.divider }} /></td>
                      <td className="px-6 py-4"><Skeleton className="h-4 w-16" style={{ backgroundColor: colors.divider }} /></td>
                      <td className="px-6 py-4"><Skeleton className="h-4 w-24" style={{ backgroundColor: colors.divider }} /></td>
                      <td className="px-6 py-4"><Skeleton className="h-4 w-24" style={{ backgroundColor: colors.divider }} /></td>
                      <td className="px-6 py-4"><Skeleton className="h-4 w-24" style={{ backgroundColor: colors.divider }} /></td>
                    </tr>
                  ))
                ) : recovered && recovered.length > 0 ? (
                  recovered.map((r, i) => (
                    <tr key={i} className="border-b last:border-0 transition-colors" style={{ borderColor: colors.divider }} onMouseEnter={(e) => e.currentTarget.style.backgroundColor = colors.hover} onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}>
                      <td className="px-6 py-3.5">
                        <div className="font-medium" style={{ color: colors.textMain }}>{r.customer_name?.String || 'Unknown'}</div>
                        <div className="text-[12px] mt-0.5" style={{ color: colors.textSecondary }}>{r.customer_email?.String || r.subscription_id}</div>
                      </td>
                      <td className="px-6 py-3.5">
                        <span style={{ color: colors.textSecondary }}>{r.customer_tier?.Valid ? r.customer_tier.String : '-'}</span>
                      </td>
                      <td className="px-6 py-3.5 font-semibold" style={{ color: colors.kpi.recovered }}>
                        {new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format((r.amount_recovered?.Int64 || 0) / 100)}
                      </td>
                      <td className="px-6 py-3.5 capitalize text-[13px]" style={{ color: colors.textSecondary }}>
                        {r.recovery_channel?.Valid ? r.recovery_channel.String.replace('send_', '') : 'Manual'}
                      </td>
                      <td className="px-6 py-3.5 text-[13px]" style={{ color: colors.textAxisTable }}>
                        {r.recovered_at?.Valid ? new Date(r.recovered_at.Time).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) : new Date(r.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} className="px-6 py-12 text-center" style={{ color: colors.textSecondary }}>
                      <ShieldAlert className="w-8 h-8 mx-auto mb-3 opacity-30" />
                      No recent recoveries to display.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};
