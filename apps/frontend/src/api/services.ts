import { apiClient } from './client';
import type { ListCasesResponse, CaseDetailsResponse, MetricsResponse, Customer, CommunicationHistory } from './types';

export const getMetrics = async (): Promise<MetricsResponse> => {
  const { data } = await apiClient.get('/metrics');
  return data;
};

export const getCases = async (page = 1, limit = 20): Promise<ListCasesResponse> => {
  const { data } = await apiClient.get('/cases', { params: { page, limit } });
  return data;
};

export const getCaseDetails = async (id: string): Promise<CaseDetailsResponse> => {
  const { data } = await apiClient.get(`/cases/${id}`);
  return data;
};

export const approveDraft = async (actionId: string) => {
  const { data } = await apiClient.post(`/cases/${actionId}/approve`);
  return data;
};

export const rejectDraft = async (actionId: string) => {
  const { data } = await apiClient.post(`/cases/${actionId}/reject`);
  return data;
};

export const getCustomer = async (id: string): Promise<Customer> => {
  const { data } = await apiClient.get(`/customers/${id}`);
  return data;
};

export const getCustomerPayments = async (id: string) => {
  const { data } = await apiClient.get(`/customers/${id}/payments`);
  return data;
};

export const getCustomerCommunications = async (id: string): Promise<CommunicationHistory[]> => {
  const { data } = await apiClient.get(`/customers/${id}/communications`);
  return data;
};
