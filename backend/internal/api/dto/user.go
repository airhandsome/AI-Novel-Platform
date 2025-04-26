package dto

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=32"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  struct {
        ID       uint   `json:"id"`
        Username string `json:"username"`
        Email    string `json:"email"`
    } `json:"user"`
}

type ResetPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type ResetPasswordConfirmRequest struct {
    Token    string `json:"token" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty" binding:"omitempty,email"`
}

type UpdatePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=6"`
} 