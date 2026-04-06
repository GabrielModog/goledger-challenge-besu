package models

type SetValueRequest struct {
	Value string `json:"value"`
}

type SetValueResponse struct {
	TxHash string `json:"tx_hash"`
	Value  string `json:"value"`
}

type GetValueResponse struct {
	Value string `json:"value"`
}

type SyncResponse struct {
	SyncedValue string `json:"synced_value"`
}

type CheckResponse struct {
	Equal      bool   `json:"equal"`
	Blockchain string `json:"blockchain_value"`
	Database   string `json:"database_value"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
