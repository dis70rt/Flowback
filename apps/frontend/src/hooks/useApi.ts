import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMetricsOverview, getMetrics, getCases, getCaseDetails, getCustomer, getCustomers, approveDraft, rejectDraft, getMetricsTrends, getMetricsChannels, getMetricsPipeline, getMetricsRecovered } from '../api/services';

export const useCustomers = () => {
  return useQuery({
    queryKey: ['customers'],
    queryFn: getCustomers,
  });
};

export const useMetrics = () => {
  return useQuery({
    queryKey: ['metrics'],
    queryFn: getMetrics,
  });
};

export const useCases = (page = 1, limit = 20) => {
  return useQuery({
    queryKey: ['cases', page, limit],
    queryFn: () => getCases(page, limit),
  });
};

export const useCaseDetails = (id?: string) => {
  return useQuery({
    queryKey: ['caseDetails', id],
    queryFn: () => getCaseDetails(id!),
    enabled: !!id,
  });
};

export const useCustomer = (id?: string) => {
  return useQuery({
    queryKey: ['customer', id],
    queryFn: () => getCustomer(id!),
    enabled: !!id,
  });
};

export const useApproveAction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: approveDraft,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['caseDetails'] });
    },
  });
};

export const useRejectAction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: rejectDraft,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['caseDetails'] });
    },
  });
};

export const useCustomerPayments = (id?: string) => {
  return useQuery({
    queryKey: ['customerPayments', id],
    queryFn: () => import('../api/services').then(m => m.getCustomerPayments(id!)),
    enabled: !!id,
  });
};

export const useCustomerCommunications = (id?: string) => {
  return useQuery({
    queryKey: ['customerCommunications', id],
    queryFn: () => import('../api/services').then(m => m.getCustomerCommunications(id!)),
    enabled: !!id,
  });
};

export const useMetricsTrends = () => {
  return useQuery({
    queryKey: ['metricsTrends'],
    queryFn: getMetricsTrends,
  });
};

export const useMetricsChannels = () => {
  return useQuery({
    queryKey: ['metricsChannels'],
    queryFn: getMetricsChannels,
  });
};

export const useMetricsPipeline = () => {
  return useQuery({
    queryKey: ['metricsPipeline'],
    queryFn: getMetricsPipeline,
  });
};

export const useMetricsRecovered = () => {
  return useQuery({
    queryKey: ['metricsRecovered'],
    queryFn: getMetricsRecovered,
  });
};

export const useMetricsOverview = () => {
  return useQuery({
    queryKey: ['metricsOverview'],
    queryFn: getMetricsOverview,
  });
};
