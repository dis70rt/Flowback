export interface NullString {
  String: string;
  Valid: boolean;
}

export interface NullInt64 {
  Int64: number;
  Valid: boolean;
}

export interface NullTime {
  Time: string;
  Valid: boolean;
}

export interface CaseItemDTO {
  id: string;
  subscription_id: string;
  amount_at_risk: number;
  status: string;
  created_at: string;
  latest_action_type: NullString;
  latest_action_status: NullString;
  latest_action_channel: NullString;
}

export interface ListCasesResponse {
  page: number;
  limit: number;
  data: CaseItemDTO[];
}

export interface MetricsResponse {
  total_revenue_recovered: number;
  active_cases: number;
  ai_success_rate: number;
}

export interface Customer {
  id: string;
  razorpay_customer_id: NullString;
  email: NullString;
  phone: NullString;
  name: NullString;
  value_tier: NullString;
  tenure: NullString;
  preferred_channel: NullString;
  total_payments: number;
  successful_payments: number;
  failed_payments: number;
  created_at: string;
  updated_at: string;
  city: NullString;
  state: NullString;
  reliability_score: number;
}

export interface NullRawMessage {
  RawMessage: {
    body?: string;
    subject?: string;
  };
  Valid: boolean;
}

export interface RecoveryAction {
  id: string;
  recovery_case_id: string;
  action_type: string;
  channel: NullString;
  ai_reasoning: NullString;
  status: string;
  draft_payload: NullRawMessage;
  created_at: string;
}

export interface RecoveryCase {
  id: string;
  customer_id: string;
  subscription_id: string;
  amount_at_risk: number;
  status: string;
  created_at: string;
}

export interface CaseDetailsResponse {
  case: RecoveryCase;
  actions: RecoveryAction[];
}

export interface CommunicationHistory {
  id: string;
  recovery_case_id: string;
  customer_id: string;
  channel: string;
  status: string;
  message_sid: NullString;
  sent_at: string;
  delivered_at: NullTime;
  opened_at: NullTime;
  clicked_at: NullTime;
}

export interface TrendData {
  date: string;
  daily_failed: number;
  daily_recovered: number;
}

export interface ChannelData {
  channel: NullString;
  count: number;
}

export interface PipelineData {
  status: string;
  count: number;
}

export interface RecoveredCase {
  id: string;
  subscription_id: string;
  payment_id: NullString;
  amount_at_risk: number;
  amount_recovered: NullInt64;
  currency: string;
  recovered_at: NullTime;
  created_at: string;
  customer_name: NullString;
  customer_email: NullString;
  customer_tier: NullString;
  recovery_channel: NullString;
  recovery_action_type: NullString;
}
