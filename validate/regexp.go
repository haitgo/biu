package validate

import (
	"regexp"
	"sync"
)

//正则式和编译集合读写锁
var regexpLock = new(sync.Mutex)

//正则式编译集合
var compiles = make(map[string]*regexp.Regexp)

//正则式集合，可以自行添加
var pattens = map[string]interface{}{
	"Any":            `^[\S\s]+$`,                                 //任意类型
	"HexCmd":         `^(\d{1,3}\t)+\d{1,3}$`,                     //16进制modbus命令
	"Str":            `^[\S]+$`,                                   //字符串
	"Date":           `^\d{4}(\-|\/|\.)\d{1,2}(\-|\/|\.)\d{1,2}$`, //这里只是一个简单的日期判断
	"Phone":          `^[1]\d{10}$`,                               //简单判断手机号码
	"Md5":            `^[a-zA-Z0-9]{32}$`,                         //md5加密
	"Email":          "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$",
	"CreditCard":     "^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$",
	"ISBN10":         "^(?:[0-9]{9}X|[0-9]{10})$",
	"ISBN13":         "^(?:[0-9]{13})$",
	"UUID3":          "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"UUID4":          "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"UUID5":          "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"UUID":           "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"Alpha":          "^[a-zA-Z]+$",
	"Alphanumeric":   "^[a-zA-Z0-9]+$",
	"Numeric":        "^[-+]?[0-9]+$",
	"Int":            "^(?:[-+]?(?:0|[1-9][0-9]*))$",
	"Float":          "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$",
	"Hexadecimal":    "^[0-9a-fA-F]+$",
	"Hexcolor":       "^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$",
	"RGBcolor":       "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$",
	"ASCII":          "^[\x00-\x7F]+$",
	"Multibyte":      "[^\x00-\x7F]",
	"FullWidth":      "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"HalfWidth":      "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"Base64":         "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$",
	"PrintableASCII": "^[\x20-\x7E]+$",
	"DataURI":        "^data:.+\\/(.+);base64$",
	"Latitude":       "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$",
	"Longitude":      "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$",
	"DNSName":        `^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{1,62}){1}(.[a-zA-Z0-9]{1}[a-zA-Z0-9_-]{1,62})*$`,
	"URL":            `^((ftp|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`,
	"SSN":            `^\d{3}[- ]?\d{2}[- ]?\d{4}$`,
	"WinPath":        `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
	"UnixPath":       `^((?:\/[a-zA-Z0-9\.\:]+(?:_[a-zA-Z0-9\:\.]+)*(?:\-[\:a-zA-Z0-9\.]+)*)+\/?)$`,
	"Semver":         "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$",
}

//注册全局验证方法
//name=用于识别的正则式标识
//patten=验证模式（可以为正则式、回调函数func(interface{})bool
/**
例：
	1、回调验证
	Reg("Name",func(name string)bool{
		return name=="wang"
	})
	2、正则式验证
	Reg("Password",`^[\d]{5}$`)
*/
func Reg(name string, patten interface{}) {
	regexpLock.Lock()
	defer regexpLock.Unlock()
	pattens[name] = patten
}

//验证
//如果在验证集合里存在该name，则调用相应的正则式，否则将patten作为正则式进行匹配
/**
例：
	1、判断是否为日期类型
	Match("12312132123","Date");
	2、注册一个验证name的回调方法
	Reg("Name",func(name string)bool{
		return name=="wang"
	})
	Match("zhang","Name");
	3、自定义正则验证
	Match("32322",`^[\d]{5}$`);
*/
func Match(value, patten string) bool {
	regexpLock.Lock()
	defer regexpLock.Unlock()
	var pat *regexp.Regexp
	var call CallFunc
	if p := compiles[patten]; pat != nil {
		pat = p
	} else {
		if p := pattens[patten]; p != nil {
			switch p.(type) {
			case string:
				pat, _ = regexp.Compile(p.(string))
			case CallFunc:
				call = p.(CallFunc)
			}
		} else {
			pat, _ = regexp.Compile(patten)
		}
	}
	if pat != nil {
		return pat.Match([]byte(value))
	} else if call != nil {
		return call(value)
	}
	return false
}
