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

export interface RecoveryAction {
  id: string;
  recovery_case_id: string;
  action_type: string;
  channel: NullString;
  status: string;
  draft_subject: NullString;
  draft_body: NullString;
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
