package clickhouse

import (
	"errors"
	"testing"

	"github.com/Parnishkaspb/curs-abds/internal/kafka"
)

// ===== Mock =====

type mockClickHouseRepo struct {
	addAcceptedCalled bool
	addDeclineCalled  bool

	addAcceptedErr error
	addDeclineErr  error
}

func (m *mockClickHouseRepo) AddAcceptedTransaction(req kafka.TransactionRequest) error {
	m.addAcceptedCalled = true
	return m.addAcceptedErr
}

func (m *mockClickHouseRepo) AddDeclineTransaction(req kafka.TransactionRequest) error {
	m.addDeclineCalled = true
	return m.addDeclineErr
}

// ===== Tests =====

func TestClickHouseService_AddAcceptedTransaction(t *testing.T) {
	tests := []struct {
		name          string
		repoErr       error
		wantErr       bool
		wantRepoCall  bool
		wantDeclineNo bool
	}{
		{
			name:          "success",
			repoErr:       nil,
			wantErr:       false,
			wantRepoCall:  true,
			wantDeclineNo: true,
		},
		{
			name:          "repo returns error",
			repoErr:       errors.New("repo error"),
			wantErr:       true,
			wantRepoCall:  true,
			wantDeclineNo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockClickHouseRepo{
				addAcceptedErr: tt.repoErr,
			}
			svc := NewClickService(repo)

			req := kafka.TransactionRequest{}

			err := svc.AddAcceptedTransaction(req)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.addAcceptedCalled != tt.wantRepoCall {
				t.Fatalf("expected AddAcceptedTransaction called=%v, got %v",
					tt.wantRepoCall, repo.addAcceptedCalled)
			}

			if tt.wantDeclineNo && repo.addDeclineCalled {
				t.Fatalf("AddDeclineTransaction should not be called")
			}
		})
	}
}

func TestClickHouseService_AddDeclineTransaction(t *testing.T) {
	tests := []struct {
		name         string
		repoErr      error
		wantErr      bool
		wantRepoCall bool
		wantAcceptNo bool
	}{
		{
			name:         "success",
			repoErr:      nil,
			wantErr:      false,
			wantRepoCall: true,
			wantAcceptNo: true,
		},
		{
			name:         "repo returns error",
			repoErr:      errors.New("repo error"),
			wantErr:      true,
			wantRepoCall: true,
			wantAcceptNo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockClickHouseRepo{
				addDeclineErr: tt.repoErr,
			}
			svc := NewClickService(repo)

			req := kafka.TransactionRequest{}

			err := svc.AddDeclineTransaction(req)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.addDeclineCalled != tt.wantRepoCall {
				t.Fatalf("expected AddDeclineTransaction called=%v, got %v",
					tt.wantRepoCall, repo.addDeclineCalled)
			}

			if tt.wantAcceptNo && repo.addAcceptedCalled {
				t.Fatalf("AddAcceptedTransaction should not be called")
			}
		})
	}
}
