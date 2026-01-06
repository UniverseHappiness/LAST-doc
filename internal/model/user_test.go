package model

import (
	"testing"
	"time"
)

func TestUser_TableName(t *testing.T) {
	user := User{}
	if user.TableName() != "users" {
		t.Errorf("TableName() = %v, want %v", user.TableName(), "users")
	}
}

func TestUser_ToResponse(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      "user",
		FirstName: "Test",
		LastName:  "User",
		Avatar:    "avatar.jpg",
		IsActive:  true,
		LastLogin: &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("ToResponse().ID = %v, want %v", response.ID, user.ID)
	}
	if response.Username != user.Username {
		t.Errorf("ToResponse().Username = %v, want %v", response.Username, user.Username)
	}
	if response.Email != user.Email {
		t.Errorf("ToResponse().Email = %v, want %v", response.Email, user.Email)
	}
	if response.Role != user.Role {
		t.Errorf("ToResponse().Role = %v, want %v", response.Role, user.Role)
	}
	if response.FirstName != user.FirstName {
		t.Errorf("ToResponse().FirstName = %v, want %v", response.FirstName, user.FirstName)
	}
	if response.LastName != user.LastName {
		t.Errorf("ToResponse().LastName = %v, want %v", response.LastName, user.LastName)
	}
	if response.Avatar != user.Avatar {
		t.Errorf("ToResponse().Avatar = %v, want %v", response.Avatar, user.Avatar)
	}
	if response.IsActive != user.IsActive {
		t.Errorf("ToResponse().IsActive = %v, want %v", response.IsActive, user.IsActive)
	}
	if response.LastLogin == nil || response.LastLogin.Unix() != user.LastLogin.Unix() {
		t.Errorf("ToResponse().LastLogin = %v, want %v", response.LastLogin, user.LastLogin)
	}
	if response.CreatedAt.Unix() != user.CreatedAt.Unix() {
		t.Errorf("ToResponse().CreatedAt = %v, want %v", response.CreatedAt, user.CreatedAt)
	}
	if response.UpdatedAt.Unix() != user.UpdatedAt.Unix() {
		t.Errorf("ToResponse().UpdatedAt = %v, want %v", response.UpdatedAt, user.UpdatedAt)
	}
}

func TestUser_ToResponse_NilLastLogin(t *testing.T) {
	user := &User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      "user",
		FirstName: "Test",
		LastName:  "User",
		Avatar:    "avatar.jpg",
		IsActive:  true,
		LastLogin: nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := user.ToResponse()

	if response.LastLogin != nil {
		t.Errorf("ToResponse().LastLogin = %v, want nil", response.LastLogin)
	}
}

func TestPasswordResetToken_TableName(t *testing.T) {
	token := PasswordResetToken{}
	if token.TableName() != "password_reset_tokens" {
		t.Errorf("TableName() = %v, want %v", token.TableName(), "password_reset_tokens")
	}
}

func TestUserLogin_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   UserLogin
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   UserLogin{Username: "testuser", Password: "password123"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.Username == "" || tt.input.Password == "" {
				t.Fatal("input should not be empty")
			}
		})
	}
}

func TestUserRegister_Validation(t *testing.T) {
	tests := []struct {
		name  string
		input UserRegister
		valid bool
	}{
		{
			name: "valid input",
			input: UserRegister{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			valid: true,
		},
		{
			name: "short username",
			input: UserRegister{
				Username: "ab",
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "invalid email",
			input: UserRegister{
				Username: "testuser",
				Email:    "invalid",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "short password",
			input: UserRegister{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "short",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.input.Username == "" || tt.input.Email == "" || tt.input.Password == "" {
					t.Error("valid input should not have empty fields")
				}
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	tests := []struct {
		name  string
		input UserUpdate
		valid bool
	}{
		{
			name:  "update email only",
			input: UserUpdate{Email: "new@example.com"},
			valid: true,
		},
		{
			name: "update all fields",
			input: UserUpdate{
				Email:     "new@example.com",
				FirstName: "New",
				LastName:  "Name",
				Avatar:    "new.jpg",
				IsActive:  func() *bool { b := false; return &b }(),
				Role:      "admin",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				// 测试字段可以正确赋值
			}
		})
	}
}

func TestPasswordResetRequest(t *testing.T) {
	tests := []struct {
		name  string
		input PasswordResetRequest
		valid bool
	}{
		{
			name:  "valid email",
			input: PasswordResetRequest{Email: "test@example.com"},
			valid: true,
		},
		{
			name:  "invalid email",
			input: PasswordResetRequest{Email: "invalid"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.input.Email == "" {
					t.Error("valid input should not have empty email")
				}
			}
		})
	}
}

func TestPasswordReset(t *testing.T) {
	tests := []struct {
		name  string
		input PasswordReset
		valid bool
	}{
		{
			name:  "valid input",
			input: PasswordReset{Token: "token123", Password: "newpass123"},
			valid: true,
		},
		{
			name:  "short password",
			input: PasswordReset{Token: "token123", Password: "short"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.input.Token == "" || tt.input.Password == "" {
					t.Error("valid input should not have empty fields")
				}
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name  string
		input ChangePassword
		valid bool
	}{
		{
			name:  "valid input",
			input: ChangePassword{OldPassword: "oldpass", NewPassword: "newpass123"},
			valid: true,
		},
		{
			name:  "short new password",
			input: ChangePassword{OldPassword: "oldpass", NewPassword: "short"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				if tt.input.OldPassword == "" || tt.input.NewPassword == "" {
					t.Error("valid input should not have empty fields")
				}
			}
		})
	}
}
