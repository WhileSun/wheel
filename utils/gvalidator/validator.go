package gvalidator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator // 全局验证器

// 单独使用
func Validate(st interface{}) error {
	uni := ut.New(zh.New()) // 翻译 实例化
	trans, _ := uni.GetTranslator("zh")

	validate := validator.New() //validator 实例化
	//err := validate.RegisterValidation("checkName", checkName)  //注册自定义函数
	err := zhTrans.RegisterDefaultTranslations(validate, trans) // 注册默认的翻译方法
	if err != nil {
		panic("gvalidator error:" + err.Error())
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string { //解析字符 标签映射  提供的方法
		label := field.Tag.Get("label")
		return label
	})
	err = validate.Struct(st)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return errors.New(v.Translate(trans))
		}
	}
	return nil
}

// 初始化gin 验证器
func InitGinValidator(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//更改TagName binding->validate
		v.SetTagName("validate")
		//v.RegisterValidation("topicurl", TopicUrl)
		v.RegisterTagNameFunc(func(field reflect.StructField) string { //解析字符 标签映射  提供的方法
			label := field.Tag.Get("label")
			return label
		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		uni := ut.New(enT, zhT, enT)
		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		switch locale {
		case "en":
			err = enTrans.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTrans.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTrans.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

// GinValidate 参数验证
func GinValidate(req *gin.Context, st interface{}) (err error) {
	// 访问格式校验
	if err = req.ShouldBind(st); err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return errors.New(v.Translate(Trans))
		}
		return
	}
	return
}
