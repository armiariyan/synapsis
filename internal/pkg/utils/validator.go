package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var v *validator.Validate

func init() {
	v = validator.New()
}

func Validate(c echo.Context, s interface{}) (err error) {
	ctx := c.Request().Context()

	if err = c.Bind(s); err != nil {
		log.Error(ctx, "error bind", err.Error())
		err = fmt.Errorf("%s", "Something Went Wrong")
		return
	}

	errVal := c.Validate(s)
	if errVal != nil {
		err = castedValidate(errVal.(validator.ValidationErrors))
		log.Error(ctx, "error validate [a]", err.Error())
		c.Set("invalid-format", true)
		return
	}

	return
}

func castedValidate(valErr validator.ValidationErrors) (err error) {
	for _, v := range valErr {
		switch v.Tag() {
		case "required":
			err = fmt.Errorf("%s is required", v.Field())
		case "email":
			err = fmt.Errorf("%s is not a valid email", v.Field())
		case "min":
			err = fmt.Errorf("%s is too short, minimum %s digit", v.Field(), v.Param())
		case "max":
			err = fmt.Errorf("%s is too long maximum %s digit", v.Field(), v.Param())
		case "len":
			err = fmt.Errorf("%s length is not valid", v.Field())
		case "eqfield":
			err = fmt.Errorf("%s is not equal to %s", v.Field(), v.Param())
		case "eq":
			err = fmt.Errorf("%s is not equal to %s", v.Field(), v.Param())
		case "gt":
			err = fmt.Errorf("%s is not greater than %s", v.Field(), v.Param())
		case "gte":
			err = fmt.Errorf("%s is not greater than or equal to %s", v.Field(), v.Param())
		case "lt":
			err = fmt.Errorf("%s is not less than %s", v.Field(), v.Param())
		case "lte":
			err = fmt.Errorf("%s is not less than or equal to %s", v.Field(), v.Param())
		case "ne":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "nfeq":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "oneof":
			err = fmt.Errorf("%s is not one of %s", v.Field(), v.Param())
		case "uuid":
			err = fmt.Errorf("%s is not a valid uuid", v.Field())
		case "ISO8601Date":
			err = fmt.Errorf("%s is not a valid ISO8601Date", v.Field())
		case "nefield":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "validInprogressStatus":
			err = fmt.Errorf("field is not valid status")
		default:
			err = fmt.Errorf("%s is not valid", v.Field())
		}
	}
	return
}

func ValidateUUID(ctx context.Context, u string) (err error) {
	_, err = uuid.Parse(u)
	if err != nil {
		log.Error(ctx, "failed to parse uuid", err)
		err = fmt.Errorf("invalid uuid")
		return
	}
	return
}

func ValidatePhoneNumberStartWith62(phoneNumber string) bool {
	matched, _ := regexp.MatchString(`^62\d+`, phoneNumber)
	return matched
}

func ValidatePasswordsEquals(password, confirmPassword string) bool {
	return password != confirmPassword
}
