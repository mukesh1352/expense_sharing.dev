export type SplitType = "EQUAL" | "EXACT" | "PERCENT";

export interface UserView {
  id: string;
  name: string;
}

export interface GroupView {
  id: string;
  name: string;
}

export interface BalanceView {
  from_user_id: string;
  to_user_id: string;
  amount: number;
}

export interface SplitInput {
  user_id: string;
  amount?: number;
  percentage?: number;
}

export type ExpenseInput = {
  expense_id: string;
  group_id: string;
  paid_by: string;
  total_amount: number;
  split_type: "EQUAL" | "EXACT" | "PERCENT";
  participants: string[];
  splits: SplitInput[]; 
};

export interface SettlementInput {
  from_user_id: string;
  to_user_id: string;
  amount: number;
}
