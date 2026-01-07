package service_test

import (
	"context"
	"errors"
	service "github.com/Parnishkaspb/curs-abds/internal/service/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// ---------- MOCK REPOSITORY ----------
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SetCountry(ctx context.Context, key, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockRepo) GetCountry(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// ---------- TESTS ----------

func TestSaveLastCountry(t *testing.T) {
	type args struct {
		accountID uint64
		country   string
		ttl       time.Duration
	}

	tests := []struct {
		name       string
		args       args
		mock       func(m *MockRepo)
		wantErr    bool
		errMessage string
	}{
		{
			name:       "invalid country",
			args:       args{accountID: 10, country: "", ttl: time.Minute},
			mock:       func(m *MockRepo) {},
			wantErr:    true,
			errMessage: "invalid county",
		},
		{
			name: "success write",
			args: args{accountID: 123, country: "RU", ttl: time.Minute},
			mock: func(m *MockRepo) {
				m.On("SetCountry",
					mock.Anything,
					"fraud:country:123",
					"RU",
					time.Minute,
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo returns error",
			args: args{accountID: 123, country: "RU", ttl: time.Minute},
			mock: func(m *MockRepo) {
				m.On("SetCountry",
					mock.Anything,
					"fraud:country:123",
					"RU",
					time.Minute,
				).Return(errors.New("redis error"))
			},
			wantErr:    true,
			errMessage: "redis error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockRepo)
			s := service.NewCountryService(m)

			tt.mock(m)

			err := s.SaveLastCountry(context.Background(),
				tt.args.accountID,
				tt.args.country,
				tt.args.ttl,
			)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestGetLastCountry(t *testing.T) {
	tests := []struct {
		name       string
		accountID  uint64
		mock       func(m *MockRepo)
		wantValue  string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "country not found (redis.Nil)",
			accountID: 123,
			mock: func(m *MockRepo) {
				m.On("GetCountry",
					mock.Anything,
					"fraud:country:123",
				).Return("", redis.Nil)
			},
			wantErr:    true,
			errMessage: service.ErrCountryNotFound.Error(),
		},
		{
			name:      "repo error",
			accountID: 123,
			mock: func(m *MockRepo) {
				m.On("GetCountry",
					mock.Anything,
					"fraud:country:123",
				).Return("", errors.New("redis error"))
			},
			wantErr:    true,
			errMessage: "redis error",
		},
		{
			name:      "success",
			accountID: 123,
			mock: func(m *MockRepo) {
				m.On("GetCountry",
					mock.Anything,
					"fraud:country:123",
				).Return("US", nil)
			},
			wantValue: "US",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockRepo)
			s := service.NewCountryService(m)

			tt.mock(m)

			val, err := s.GetLastCountry(context.Background(), tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, val)
			}

			m.AssertExpectations(t)
		})
	}
}
