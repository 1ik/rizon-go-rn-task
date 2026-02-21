// Only EmailAuth is called. Table-driven: one loop, expect error or token.
package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"rizon-test-task/internal/config"
	"rizon-test-task/internal/mocks"
	"rizon-test-task/internal/models"
	"rizon-test-task/internal/repository"

	"go.uber.org/mock/gomock"
)

var errAny = errors.New("any error")

func TestEmailAuth(t *testing.T) {
	const testEmail = "test@example.com"
	const storedHash = "stored-hash-secret"
	ctx := context.Background()
	authCfg := &config.AuthConfig{JWTSecret: "test-secret", JWTExpiration: time.Hour}

	tests := map[string]struct {
		email         string
		secret        string
		setupStore    func(*mocks.MockStore)
		setupUserRepo func(*mocks.MockUserRepository)
		wantErr       error // nil = success
		wantErrMsg    string
		wantToken     bool
	}{
		"invalid email format": {
			email:  "not-an-email",
			secret: "any",
			setupStore: func(s *mocks.MockStore) {},
			setupUserRepo: func(r *mocks.MockUserRepository) {},
			wantErr:    errAny,
			wantErrMsg: "invalid email format",
		},
		"store Exists error": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(false, errors.New("store error"))
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {},
			wantErr:       errAny,
		},
		"email not in store": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(false, nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {},
			wantErr:       ErrEmailAuthNotFound,
		},
		"store Get error": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return("", errors.New("get error"))
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {},
			wantErr:       errAny,
		},
		"secret mismatch": {
			email:  testEmail,
			secret: "wrong-secret",
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {},
			wantErr:       ErrEmailAuthInvalidSecret,
		},
		"new user Create success": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
				s.EXPECT().Delete(gomock.Any(), testEmail).Return(nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByEmail(gomock.Any(), testEmail).Return(nil, repository.ErrUserNotFound)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, u *models.User) error {
					u.ID = 1
					return nil
				})
			},
			wantToken: true,
		},
		"existing user": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
				s.EXPECT().Delete(gomock.Any(), testEmail).Return(nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByEmail(gomock.Any(), testEmail).Return(&models.User{ID: 42, Email: testEmail}, nil)
			},
			wantToken: true,
		},
		"FindByEmail error": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByEmail(gomock.Any(), testEmail).Return(nil, errors.New("db error"))
			},
			wantErr: errAny,
		},
		"Create fails": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByEmail(gomock.Any(), testEmail).Return(nil, repository.ErrUserNotFound)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("create failed"))
			},
			wantErr: errAny,
		},
		"Store Delete fails still returns token": {
			email:  testEmail,
			secret: storedHash,
			setupStore: func(s *mocks.MockStore) {
				s.EXPECT().Exists(gomock.Any(), testEmail).Return(true, nil)
				s.EXPECT().Get(gomock.Any(), testEmail).Return(storedHash, nil)
				s.EXPECT().Delete(gomock.Any(), testEmail).Return(errors.New("delete failed"))
			},
			setupUserRepo: func(r *mocks.MockUserRepository) {
				r.EXPECT().FindByEmail(gomock.Any(), testEmail).Return(&models.User{ID: 1, Email: testEmail}, nil)
			},
			wantToken: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocks.NewMockStore(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			tc.setupStore(store)
			tc.setupUserRepo(userRepo)
			a := &appImpl{
				userRepo:      userRepo,
				feedbackRepo:  mocks.NewMockFeedbackRepository(ctrl),
				store:         store,
				authCfg:       authCfg,
				messageBroker: mocks.NewMockMessageBroker(ctrl),
			}
			token, err := a.EmailAuth(ctx, tc.email, tc.secret)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatal("EmailAuth() error = nil, want error")
				}
				if tc.wantErr != errAny && !errors.Is(err, tc.wantErr) {
					t.Errorf("EmailAuth() error = %v, want %v", err, tc.wantErr)
				}
				if tc.wantErrMsg != "" && err.Error() != tc.wantErrMsg {
					t.Errorf("EmailAuth() error = %v, want msg %q", err, tc.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("EmailAuth() unexpected error = %v", err)
			}
			if tc.wantToken && token == "" {
				t.Fatal("EmailAuth() token is empty")
			}
		})
	}
}
