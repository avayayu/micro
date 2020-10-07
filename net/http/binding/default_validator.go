package binding

import (
	"reflect"
	"strings"
	"sync"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器T
var trans ut.Translator

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) RegisterValidation(key string, fn validator.Func) error {
	v.lazyinit()
	return v.validate.RegisterValidation(key, fn)
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {

		zh := zhongwen.New()
		uni := ut.New(zh, zh)
		trans, _ = uni.GetTranslator("zh")

		validate := validator.New()
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			label := field.Tag.Get("label")
			if label == "" {
				return field.Name
			}
			return label
		})
		zh_translations.RegisterDefaultTranslations(validate, trans)
		v.validate = validate

	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *defaultValidator) GetValidate() *validator.Validate {
	v.lazyinit()

	return v.validate
}

func Translate(err error) string {
	var errs validator.ValidationErrors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return "无法识别错误"
	}

	var errList []string
	for _, e := range errs {

		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}
