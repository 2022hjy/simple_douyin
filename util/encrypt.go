package util

import (
    "golang.org/x/crypto/bcrypt"
)

// EncryptPassword 将密码加密，需要传入密码返回的是加密后的密码
func EncryptPassword(password string) (string, error) {
    // 加密密码，使用 bcrypt 包当中的 GenerateFromPassword 方法，bcrypt.DefaultCost 代表使用默认加密成本
    encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    if err != nil {
        // 如果有错误则返回异常，加密后的空字符串返回为空字符串，因为加密失败
        return "", err
    }
    // 返回加密后的密码和空异常
    return string(encryptPassword), nil
}

// ValidatePassword 验证密码，需要传入加密后的密码和密码，返回的是是否验证成功
func ValidatePassword(encryptPassword string, password string) bool {
    // 使用 bcrypt 包当中的 CompareHashAndPassword 方法来验证密码是否正确
    err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password))
    return err == nil
}
