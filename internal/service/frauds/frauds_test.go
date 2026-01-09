package frauds

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
	"github.com/Parnishkaspb/curs-abds/internal/service/clickhouse"
)

type mockCHRepo struct {
	acceptedCalled int
	declineCalled  int

	lastAccepted kafka.TransactionRequest
	lastDecline  kafka.TransactionRequest

	acceptedErr error
	declineErr  error
}

func (m *mockCHRepo) AddAcceptedTransaction(req kafka.TransactionRequest) error {
	m.acceptedCalled++
	m.lastAccepted = req
	return m.acceptedErr
}

func (m *mockCHRepo) AddDeclineTransaction(req kafka.TransactionRequest) error {
	m.declineCalled++
	m.lastDecline = req
	return m.declineErr
}

// helper to marshal request
func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return string(b)
}

func TestFrauds_checkHighAmount(t *testing.T) {
	f := &Frauds{}

	tests := []struct {
		name       string
		amount     uint64
		wantDecl   bool
		wantReason string
		wantLog    bool
	}{
		{
			name:       "amount below threshold -> accept",
			amount:     499999,
			wantDecl:   false,
			wantReason: "",
			wantLog:    false,
		},
		{
			name:       "amount equals threshold -> decline",
			amount:     500000,
			wantDecl:   true,
			wantReason: "high_amount",
			wantLog:    true,
		},
		{
			name:       "amount above threshold -> decline",
			amount:     999999,
			wantDecl:   true,
			wantReason: "high_amount",
			wantLog:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := kafka.TransactionRequest{Amount: tt.amount}

			got := f.checkHighAmount(context.Background(), req)

			if got.Decline != tt.wantDecl {
				t.Fatalf("Decline: got %v want %v", got.Decline, tt.wantDecl)
			}
			if got.Reason != tt.wantReason {
				t.Fatalf("Reason: got %q want %q", got.Reason, tt.wantReason)
			}
			if tt.wantLog && got.Log == "" {
				t.Fatalf("expected Log to be non-empty")
			}
			if !tt.wantLog && got.Log != "" {
				t.Fatalf("expected Log to be empty, got %q", got.Log)
			}
		})
	}
}

func TestFrauds_blacklist(t *testing.T) {
	f := &Frauds{}

	tests := []struct {
		name       string
		merchant   string
		wantDecl   bool
		wantReason string
		wantLog    bool
	}{
		{
			name:       "merchant not in blacklist -> accept",
			merchant:   "OZON",
			wantDecl:   false,
			wantReason: "",
			wantLog:    false,
		},
		{
			name:       "merchant YMARKET -> decline",
			merchant:   "YMARKET",
			wantDecl:   true,
			wantReason: "blacklist",
			wantLog:    true,
		},
		{
			name:       "merchant CAMOKAT -> decline",
			merchant:   "CAMOKAT",
			wantDecl:   true,
			wantReason: "blacklist",
			wantLog:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := kafka.TransactionRequest{Merchant: tt.merchant}

			got := f.blacklist(context.Background(), req)

			if got.Decline != tt.wantDecl {
				t.Fatalf("Decline: got %v want %v", got.Decline, tt.wantDecl)
			}
			if got.Reason != tt.wantReason {
				t.Fatalf("Reason: got %q want %q", got.Reason, tt.wantReason)
			}
			if tt.wantLog && got.Log == "" {
				t.Fatalf("expected Log to be non-empty")
			}
			if !tt.wantLog && got.Log != "" {
				t.Fatalf("expected Log to be empty, got %q", got.Log)
			}
		})
	}
}

// ---------- CHECKMESSAGE / PROCESS TESTS (table-driven) ----------
func TestFrauds_CheckMessage_ClickHouseRouting(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		rules           []Rule
		wantAccepted    int
		wantDecline     int
		wantAcceptedReq bool
		wantDeclineReq  bool
	}{
		{
			name:    "bad json -> no clickhouse calls",
			message: "{not-json",
			rules: []Rule{
				{
					Code: "any",
					Func: func(_ context.Context, _ kafka.TransactionRequest) FraudResult {
						return FraudResult{Decline: false}
					},
				},
			},
			wantAccepted: 0,
			wantDecline:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chRepo := &mockCHRepo{}
			chSvc := clickhouse.NewClickService(chRepo)

			f := &Frauds{
				rules:      tt.rules,
				clickhouse: chSvc,
			}

			f.CheckMessage(tt.message)

			if chRepo.acceptedCalled != tt.wantAccepted {
				t.Fatalf("acceptedCalled: got %d want %d", chRepo.acceptedCalled, tt.wantAccepted)
			}
			if chRepo.declineCalled != tt.wantDecline {
				t.Fatalf("declineCalled: got %d want %d", chRepo.declineCalled, tt.wantDecline)
			}

			if tt.wantAcceptedReq {
				var want kafka.TransactionRequest
				_ = json.Unmarshal([]byte(tt.message), &want)
				if chRepo.lastAccepted != want {
					t.Fatalf("accepted req mismatch: got %+v want %+v", chRepo.lastAccepted, want)
				}
			}
			if tt.wantDeclineReq {
				var want kafka.TransactionRequest
				_ = json.Unmarshal([]byte(tt.message), &want)
				if chRepo.lastDecline != want {
					t.Fatalf("decline req mismatch: got %+v want %+v", chRepo.lastDecline, want)
				}
			}
		})
	}
}
