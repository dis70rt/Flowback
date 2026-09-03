import { apiClient } from './client';
import type { ListCasesResponse, CaseDetailsResponse, MetricsResponse, Customer, CommunicationHistory, TrendData, ChannelData, PipelineData, RecoveredCase } from './types';

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

export const getCustomers = async (): Promise<Customer[]> => {
  const { data } = await apiClient.get('/customers');
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

export const getMetricsTrends = async (): Promise<TrendData[]> => {
  const { data } = await apiClient.get('/metrics/trends');
  return data;
};

export const getMetricsChannels = async (): Promise<ChannelData[]> => {
  const { data } = await apiClient.get('/metrics/channels');
  return data;
};

export const getMetricsPipeline = async (): Promise<PipelineData[]> => {
  const { data } = await apiClient.get('/metrics/pipeline');
  return data;
};

export const getMetricsRecovered = async (): Promise<RecoveredCase[]> => {
  const { data } = await apiClient.get('/metrics/recovered');
  return data;
};

export interface OverviewData {
  total_amount_at_risk: number;
  total_amount_recovered: number;
  active_cases: number;
  recovered_cases: number;
  ai_success_rate: number;
}
export const getMetricsOverview = async (): Promise<OverviewData> => {
  const { data } = await apiClient.get('/metrics/overview');
  return data;
};
