package models

type Credits struct {
	CurrentBalance          int `json:"current_balance"`
	EstimatedDailyIncome    int `json:"estimated_daily_income"`
	EstimatedDailyExpense   int `json:"estimated_daily_expenditure"`
	EstimatedDailyBalance   int `json:"estimated_daily_balance"`
	CalculationTime         string `json:"calculation_time"`
	EstimatedRunoutSeconds  int   `json:"estimated_runout_seconds,omitempty"`
	PastDayMeasurementSpend int   `json:"past_day_measurement_results,omitempty"`
	PastDayTransfer         int   `json:"past_day_credits_spent,omitempty"`
	IncomeItems             string `json:"income_items,omitempty"`
	ExpenseItems            string `json:"expense_items,omitempty"`
}


type TransferRequest struct {
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}

type TransferResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
