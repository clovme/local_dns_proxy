package boot

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"local_dns_proxy/internal/bootstrap/boot/validate"
)

// InitializationFormValidate 初始化表单验证器
func InitializationFormValidate() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		_ = v.RegisterValidation("emailValid", validate.EmailValid)
		_ = v.RegisterValidation("usernameValid", validate.UsernameValid)
		_ = v.RegisterValidation("passwordValid", validate.PasswordValid)
	}
}
