import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMetrics, getCases, getCaseDetails, getCustomer, approveDraft, rejectDraft } from '../api/services';

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
