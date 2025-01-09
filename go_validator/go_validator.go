package main

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

func main() {}

/*
多个校验标签将会按照定义顺序依次校验。
校验器按逗号（and）、管道符（or）分隔，如果校验参数包含逗号，可使用 0x2c 代替，管道符 | 使用 0x7c 代替。

-：跳过校验。
structonly
nostructlevel
omitempty：如果字段值是零值，跳过校验。
omitnil：如果字段值是 nil，跳过校验。
dive：进入切片、数组和映射字段内校验。例如：gt=0,dive,keys,dive,eq=1|eq=2,endkeys,required。
keys：进入映射字段的键校验。
endkeys：结束映射字段的键校验。

required：要求字段不是零值。
required_if：如果另一个字段的值等于某个值，则字段不能为零值。例如：required_if=Field1 foo Field2 bar。
required_unless：除非另一个字段的值等于某个值，则字段不能为零值。
required_with：如果别的字段不是零值，则字段不能为零值。例如：required_with=Field1 Field2。
required_with_all：如果所有别的字段不是零值，则字段不能为零值。
required_without：如果别的字段是零值，则字段不能为零值。
required_without_all：如果所有别的字段是零值，则字段不能为零值。
excluded_if：如果所有别的字段等于某个值，则排除字段校验。
excluded_unless：如果所有别的字段等于某个值，则字段不能为零值。

notblank
isdefault：校验是否是默认值。
oneof：要求字段值是其中一个。
oneofci：同上，但忽略大小写。
unique：要求数组、切片和映射没有相同的元素。要求结构体切片的指定字段没有相同元素。
len：要求数字、字符串、数组、切片、时长和映射的长度等于某个值。
max：要求数字、字符串、数组、切片、时长和映射的长度不超过某个值。
min：与 max 语义相反。

eq：同 len 语义。
eq_ignore_case
ne：与 eq 语义相反。
ne_ignore_case
gt：数字大于，字符串、数组、切片、时间、时长和映射的长度大于。
gte：同上，大于等于。与 min 相同。
lt：与 gt 语义相反。
lte：与 gte 语义相反。

eqfield：要求与另一个字段值相同。
eqcsfield：要求与另一个可以不在同层次的字段值相同。
nefield：与 eqfield 语义相反。
necsfield：与 eqcsfield 语义相反。
gtfield：要求大于另一个数字、时长或时间类型字段的值。
gtcsfield：与 gtfield 语义相同，但允许比较不同层次的字段。
gtefield：同 gtfield 语义，允许等于。
gtecsfield：同 gtcsfield 语义，允许等于。
ltfield：与 gtfield 语义相反。
ltcsfield：与 gtcsfield 语义相反。
ltefield：与 gtefield 语义相反。
ltecsfield：与 gtecsfield 语义相反。
containsfield：要求包含另一个字段的字符串值。
excludesfield：与 containsfield 语义相反。

contains：要求字符串包含某个子串。
excludes：要求字符串不包含某个子串。
containsany：要求字符串包含任何一个字符。
containsrune：要求字符串包含某个字符。
excludesrune：要求字符串不包含某个字符。
excludesall：要求字符串包含任何一个字符。
startswith：要求字符串以某子串开头。
startsnotwith：与上语义相反。
endswith：要求字符串以某子串结尾。
endsnotwith：与上语义相反。

alpha：要求字符串只能包含 ASCII 码字符。
alphanum：要求字符串只能包含 ASCII 码中的字母数字字符。
alphaunicode：要求字符串都是 Unicode 编码字符。
alphanumunicode：要求字符串都是 Unicode 编码字符。
boolean：要求字符串是布尔值。
number：要求字符串是数字值。
numeric：要求字符串是基础数字值。不包含指数 E。
hexadecimal：要求字符串是十六进制值。
lowercase：要求字符串是小写字母，不能是空串。
uppercase：要求字符串是大写字母，不能是空串。
ascii：要求字符串都是 ASCII 字符，可为空串。
printascii：要求字符串是可打印的 ASCII 字符，可为空串。
multibyte：要求字符串包含多字节字符，可为空串。

hexcolor：要求字符串是颜色编码。包含 #。
rgb：要求字符串是 RGB 颜色编码。
rgba：要求字符串是 RGBA 颜色编码。
hsl：要求字符串是 HSL 颜色编码。
hsla：要求字符串是 HSLA 颜色编码。
iscolor：hexcolor|rgb|rgba|hsl|hsla

e164：手机号。例如：+8613812345678。
email：要求字符串是邮件地址。
json：要求字符串是 JSON 格式。
jwt：要求字符串是 JWT 格式。

file：要求字符串是文件路径，且系统上存在该文件。
image：要求字符串是图片文件路径，且系统上存在该文件。
filepath：要求字符串是文件路径。
dir：要求字符串是文件夹路径，且系统上存在该文件夹。
dirpath：要求字符串是文件夹路径。

url：要求字符串是 URL，且有协议部分。
uri：要求字符串是 URI。
datauri
http_url
html：要求字符串是 HTML 标签。
html_encoded
url_encoded：要求字符串是 URL 编码字符串。
urn_rfc2141：要求字符串是 URN。

base32：要求字符串是 Base32，不可为空串。
base64：要求字符串是 Base64，不可为空串。
base64url：要求字符串是 Base64URL，不可为空串。
base64rawurl

btc_addr：要求字符串是比特币地址。
btc_addr_bech32：要求字符串是比特币地址。
eth_addr：要求字符串是以太坊地址。

uuid：要求字符串是 UUID。
uuid_rfc4122
uuid3：同上。
uuid3_rfc4122
uuid4：同上。
uuid4_rfc4122
uuid5：同上。
uuid5_rfc4122

isbn：要求字符串是 isbn10 或 isbn13。
isbn10
isbn13
issn
ulid：Universally Unique Lexicographically Sortable Identifier
latitude：要求字符串是纬度。
longitude：要求字符串是经度。
ssn
postcode_iso3166_alpha2
postcode_iso3166_alpha2_field
fqdn：要求字符串包含全限定域名。
hostname_port
datetime：要求字符串是指定的时间模板。
bcp47_language_tag
bic：Business Identifier Code
dns_rfc1035_label
timezone：要求字符串是系统上存在的时区字符串。
semver：要求字符串是语义化版本。
cve：Common Vulnerabilities and Exposures Identifier
credit_card
luhn_checksum
mongodb
mongodb_connection_string
cron：要求字符串包含 CRON。
spicedb

md4
md5
sha256
sha384
sha512
ripemd128
ripemd128
tiger128
tiger160
tiger192
luhn_checksum

ip：要求字符串包含合法 IP 地址。
ipv4
ipv6
cidr：要求字符串包含合法 CIDR 地址。
cidrv4
cidrv6
tcp_addr：要求字符串包含可解析的 TCP 地址。
tcp4_addr
tcp6_addr
udp_addr：要求字符串包含可解析的 UDP 地址。
udp4_addr
udp6_addr
ip_addr：要求字符串包含可解析的 IP 地址。
ip4_addr
ip6_addr
unix_addr：要求字符串包含 Unix 地址。
mac
hostname
hostname_rfc1123

iso4217
iso3166_2
iso3166_1_alpha_numeric
iso3166_1_alpha2
iso3166_1_alpha3
country_code：iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric
*/

type User struct {
	Gender          string        `validate:"oneof='male' 'female' 'prefer not to'"` // 枚举值有空格，则使用单引号括起。
	Birthday        time.Time     `validate:"gt"`                                    // 要求大于 time.Now().UTC()。
	WaitTime        time.Duration `validate:"gt=1h30m"`                              // 要求时长大于。
	Password        string        `validate:"eqfield=ConfirmPassword,eqcsfield=innerInfo.Password"`
	ConfirmPassword string
	innerInfo       struct {
		Password string
	}
	Hobbies []struct {
		Name string
	} `validate:"unique=Name"`
}

func ValidateVar() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	myEmail := "ivfzhou@126.com@"
	err := vd.Var(myEmail, "required,email")
	if err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			v := errs[0]
			fmt.Println("Namespace:", v.Namespace())
			fmt.Println("StructNamespace:", v.StructNamespace())
			fmt.Println("StructField:", v.StructField())
			fmt.Println("Field:", v.Field())
			fmt.Println("Kind:", v.Kind())           // string
			fmt.Println("Tag:", v.Tag())             // email
			fmt.Println("ActualTag:", v.ActualTag()) // email
			fmt.Println("Type:", v.Type())           // string
			fmt.Println("Value:", v.Value())         // ivfzhou@126.com@
			fmt.Println("Param:", v.Param())
			fmt.Println("Error:", v.Error())
		}
	}
}

func ValidateStruct() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	type User struct {
		Name   string `validate:"required"`
		Gender string `validate:"oneof=male female prefer_not_to"`
	}
	u := &User{
		Name:   "zs",
		Gender: "man",
	}
	err := vd.Struct(u)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, v := range errs {
			fmt.Println("Namespace:", v.Namespace())             // User.Gender
			fmt.Println("StructNamespace:", v.StructNamespace()) // User.Gender
			fmt.Println("StructField:", v.StructField())         // Gender
			fmt.Println("Field:", v.Field())                     // Gender
			fmt.Println("Kind:", v.Kind())                       // string
			fmt.Println("Tag:", v.Tag())                         // oneof
			fmt.Println("ActualTag:", v.ActualTag())             // oneof
			fmt.Println("Type:", v.Type())                       // string
			fmt.Println("Value:", v.Value())                     // man
			fmt.Println("Param:", v.Param())                     // male female prefer_not_to
			fmt.Println("Error:", v.Error())
			fmt.Println()
		}
	}
}

func ValidateMap() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	m := map[string]any{
		"Name":    "",
		"Address": map[string]any{
			// "Street": "",
		},
	}
	errMap := vd.ValidateMap(m, map[string]any{
		"Name": "required",
		"Address": map[string]any{
			"Street": "required",
		},
	})
	if len(errMap) > 0 {
		for fieldName, errs := range errMap {
			fmt.Println("Error Field:", fieldName)
			switch err := errs.(type) {
			case validator.ValidationErrors: // map 外层错误
				for _, v := range err {
					fmt.Println("Namespace:", v.Namespace()) // Name
					fmt.Println("StructNamespace:", v.StructNamespace())
					fmt.Println("StructField:", v.StructField())
					fmt.Println("Field:", v.Field())
					fmt.Println("Kind:", v.Kind())           // string
					fmt.Println("Tag:", v.Tag())             // required
					fmt.Println("ActualTag:", v.ActualTag()) // required
					fmt.Println("Type:", v.Type())           // string
					fmt.Println("Value:", v.Value())
					fmt.Println("Param:", v.Param())
					fmt.Println("Error:", v.Error())
					fmt.Println()
				}
			case map[string]any: // 内层 map 错误
				for k, v := range err {
					fmt.Println("\tError Field:", k)
					errs := v.(validator.ValidationErrors)
					for _, v := range errs {
						fmt.Println("\tNamespace:", v.Namespace()) // Street
						fmt.Println("\tStructNamespace:", v.StructNamespace())
						fmt.Println("\tStructField:", v.StructField())
						fmt.Println("\tField:", v.Field())
						fmt.Println("\tKind:", v.Kind())           // invalid
						fmt.Println("\tTag:", v.Tag())             // mytag
						fmt.Println("\tActualTag:", v.ActualTag()) // required
						fmt.Println("\tType:", v.Type())
						fmt.Println("\tValue:", v.Value())
						fmt.Println("\tParam:", v.Param())
						fmt.Println("\tError:", v.Error())
						fmt.Println()
					}
				}
			}
		}
	}
}

func CustomTypeValidate() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	vd.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if valuer, ok := field.Interface().(driver.Valuer); ok {
			val, err := valuer.Value()
			if err == nil {
				// 返回要校验的值
				return val
			}
		}
		return nil
	}, sql.NullInt64{})
	s := sql.NullInt64{
		Int64: 4,
		Valid: true,
	}
	err := vd.Var(s, "required,min=1,max=3")
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())
		fmt.Println("StructNamespace:", v.StructNamespace())
		fmt.Println("StructField:", v.StructField())
		fmt.Println("Field:", v.Field())
		fmt.Println("Kind:", v.Kind())           // int64
		fmt.Println("Tag:", v.Tag())             // max
		fmt.Println("ActualTag:", v.ActualTag()) // max
		fmt.Println("Type:", v.Type())           // int64
		fmt.Println("Value:", v.Value())         // 4
		fmt.Println("Param:", v.Param())         // 3
		fmt.Println("Error:", v.Error())
	}
}

func CustomTagValidate() {
	type User struct {
		Name string `validate:"mytag"`
	}
	u := &User{
		Name: "zs",
	}
	vd := validator.New(validator.WithRequiredStructEnabled())
	// 相同 tag 的校验逻辑会被覆盖。
	err := vd.RegisterValidation("mytag", func(fl validator.FieldLevel) bool {
		fmt.Println("Field:", fl.Field().Interface()) // zs
		fmt.Println("Param:", fl.Param())
		fmt.Println("Parent:", fl.Parent().Interface())       // {zs}
		fmt.Println("FieldName:", fl.FieldName())             // Name
		fmt.Println("StructFieldName:", fl.StructFieldName()) // Name
		fmt.Println("GetTag:", fl.GetTag())                   // mytag
		fmt.Println("Top:", fl.Top().Interface())             // &{zs}
		fmt.Println()
		return fmt.Sprint(fl.Field().Interface()) == "ivfzhou"
	})
	if err != nil {
		panic(err)
	}
	err = vd.Struct(u)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())             // User.Name
		fmt.Println("StructNamespace:", v.StructNamespace()) // User.Name
		fmt.Println("StructField:", v.StructField())         // Name
		fmt.Println("Field:", v.Field())                     // Name
		fmt.Println("Kind:", v.Kind())                       // string
		fmt.Println("Tag:", v.Tag())                         // mytag
		fmt.Println("ActualTag:", v.ActualTag())             // mytag
		fmt.Println("Type:", v.Type())                       // string
		fmt.Println("Value:", v.Value())                     // zs
		fmt.Println("Param:", v.Param())
		fmt.Println("Error:", v.Error())
	}
}

func CustomStructValidate() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	type User struct {
		Name string
	}
	vd.RegisterStructValidation(func(sl validator.StructLevel) {
		user := sl.Current().Interface().(User)
		if len(user.Name) <= 0 {
			sl.ReportError(user.Name, "FieldName", "StructFieldName", "Tag", "Param")
		}
	}, User{})
	err := vd.Struct(&User{})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, v := range errs {
			fmt.Println("Namespace:", v.Namespace())             // User.FieldName
			fmt.Println("StructNamespace:", v.StructNamespace()) // User.StructFieldName
			fmt.Println("StructField:", v.StructField())         // StructFieldName
			fmt.Println("Field:", v.Field())                     // FieldName
			fmt.Println("Kind:", v.Kind())                       // string
			fmt.Println("Tag:", v.Tag())                         // Tag
			fmt.Println("ActualTag:", v.ActualTag())             // Tag
			fmt.Println("Type:", v.Type())                       // string
			fmt.Println("Value:", v.Value())
			fmt.Println("Param:", v.Param()) // Param
			fmt.Println("Error:", v.Error())
			fmt.Println()
		}
	}
}

func CustomStructValidateByMap() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	type Address struct {
		City   string
		Street string
	}
	type User struct {
		Name    string
		Address Address
	}
	vd.RegisterStructValidationMapRules(map[string]string{
		"Name":    "required",
		"Address": "required",
	}, User{})
	vd.RegisterStructValidationMapRules(map[string]string{
		"City": "required",
	}, Address{})
	err := vd.Struct(&User{
		Name: "",
		Address: Address{
			Street: "1",
		},
	})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, v := range errs {
			fmt.Println("Namespace:", v.Namespace())
			fmt.Println("StructNamespace:", v.StructNamespace())
			fmt.Println("StructField:", v.StructField())
			fmt.Println("Field:", v.Field())
			fmt.Println("Kind:", v.Kind())
			fmt.Println("Tag:", v.Tag())
			fmt.Println("ActualTag:", v.ActualTag())
			fmt.Println("Type:", v.Type())
			fmt.Println("Value:", v.Value())
			fmt.Println("Param:", v.Param())
			fmt.Println("Error:", v.Error())
			fmt.Println()
		}
	}
}

func TagAlias() {
	type User struct {
		Name string `validate:"mytag"`
	}
	u := &User{
		Name: "",
	}
	vd := validator.New(validator.WithRequiredStructEnabled())
	vd.RegisterAlias("mytag", "required")
	err := vd.Struct(u)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())             // User.Name
		fmt.Println("StructNamespace:", v.StructNamespace()) // User.Name
		fmt.Println("StructField:", v.StructField())         // Name
		fmt.Println("Field:", v.Field())                     // Name
		fmt.Println("Kind:", v.Kind())                       // string
		fmt.Println("Tag:", v.Tag())                         // mytag
		fmt.Println("ActualTag:", v.ActualTag())             // required
		fmt.Println("Type:", v.Type())                       // string
		fmt.Println("Value:", v.Value())
		fmt.Println("Param:", v.Param())
		fmt.Println("Error:", v.Error())
	}
}

func ModifyFieldName() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	vd.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	type User struct {
		Name string `json:"nickname,omitempty" validate:"required"`
	}
	u := &User{}
	err := vd.Struct(u)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())             // User.Name
		fmt.Println("StructNamespace:", v.StructNamespace()) // User.Name
		fmt.Println("StructField:", v.StructField())         // Name
		fmt.Println("Field:", v.Field())                     // nickname
		fmt.Println("Kind:", v.Kind())                       // string
		fmt.Println("Tag:", v.Tag())                         // required
		fmt.Println("ActualTag:", v.ActualTag())             // required
		fmt.Println("Type:", v.Type())                       // string
		fmt.Println("Value:", v.Value())
		fmt.Println("Param:", v.Param())
		fmt.Println("Error:", v.Error())
	}
}

func Translate() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	ts := ut.New(zh.New(), en.New())
	zht, _ := ts.GetTranslator("zh")
	err := zh_trans.RegisterDefaultTranslations(vd, zht)
	if err != nil {
		panic(err)
	}
	ent, _ := ts.GetTranslator("en")
	err = en_trans.RegisterDefaultTranslations(vd, ent)
	if err != nil {
		panic(err)
	}

	type User struct {
		Name string `validate:"required"`
	}
	err = vd.Struct(&User{})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())
		fmt.Println("StructNamespace:", v.StructNamespace())
		fmt.Println("StructField:", v.StructField())
		fmt.Println("Field:", v.Field())
		fmt.Println("Kind:", v.Kind())
		fmt.Println("Tag:", v.Tag())
		fmt.Println("ActualTag:", v.ActualTag())
		fmt.Println("Type:", v.Type())
		fmt.Println("Value:", v.Value())
		fmt.Println("Param:", v.Param())
		fmt.Println("Error:", v.Error())
		fmt.Println("zht:", v.Translate(zht)) // Name为必填字段
		fmt.Println("ent:", v.Translate(ent)) // Name is a required field
	}
}

func CustomTranslate() {
	vd := validator.New(validator.WithRequiredStructEnabled())
	ts := ut.New(zh.New(), en.New())
	zht, _ := ts.GetTranslator("zh")
	err := vd.RegisterTranslation("required", zht,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} 必须有值!", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T("required", fe.Field())
			if err != nil {
				panic(err)
			}
			return t
		},
	)
	if err != nil {
		panic(err)
	}
	type User struct {
		Name string `validate:"required"`
	}
	err = vd.Struct(&User{})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		v := errs[0]
		fmt.Println("Namespace:", v.Namespace())
		fmt.Println("StructNamespace:", v.StructNamespace())
		fmt.Println("StructField:", v.StructField())
		fmt.Println("Field:", v.Field())
		fmt.Println("Kind:", v.Kind())
		fmt.Println("Tag:", v.Tag())
		fmt.Println("ActualTag:", v.ActualTag())
		fmt.Println("Type:", v.Type())
		fmt.Println("Value:", v.Value())
		fmt.Println("Param:", v.Param())
		fmt.Println("Error:", v.Error())
		fmt.Println("zht:", v.Translate(zht)) // Name为必填字段
	}
}

type GinValidator struct {
	validator *validator.Validate
}

func (v *GinValidator) ValidateStruct(a any) error {
	return v.validator.Struct(a)
}

func (v *GinValidator) Engine() any {
	return v.validator
}

func (v *GinValidator) Init() {
	v.validator = validator.New(validator.WithRequiredStructEnabled())
	v.validator.SetTagName("binding")
}

func WithGin() {
	vd := &GinValidator{}
	vd.Init()
	binding.Validator = vd
}
