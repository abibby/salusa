package rules

import (
	"encoding/json"
	"log"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func initStringRules() {
	AddStringRule("alpha", func(value string, args []string) bool {
		for _, c := range value {
			if !alpha(c) {
				return false
			}
		}
		return true
	})
	AddStringRule("alpha_dash", func(value string, args []string) bool {
		for _, c := range value {
			if !alpha(c) && c != '-' {
				return false
			}
		}
		return true
	})
	AddStringRule("alpha_num", func(value string, args []string) bool {
		for _, c := range value {
			if !alpha(c) && !numeric(c) {
				return false
			}
		}
		return true
	})
	AddStringRule("numeric", func(value string, args []string) bool {
		for _, c := range value {
			if !numeric(c) {
				return false
			}
		}
		return true
	})
	AddStringRule("email", func(value string, args []string) bool {
		_, err := mail.ParseAddress(value)
		return err == nil
	})
	AddStringRule("ends_with", func(value string, args []string) bool {
		if len(args) < 1 {
			log.Print("end_with must have 1 argument")
			return false
		}
		return strings.HasSuffix(value, args[0])
	})
	AddStringRule("starts_with", func(value string, args []string) bool {
		if len(args) < 1 {
			log.Print("starts_with must have 1 argument")
			return false
		}
		return strings.HasPrefix(value, args[0])
	})
	AddStringRule("ip_address", func(value string, args []string) bool {
		return net.ParseIP(value) != nil
	})
	AddStringRule("json", func(value string, args []string) bool {
		var v any
		err := json.Unmarshal([]byte(value), &v)
		return err == nil
	})
	AddStringRule("mac_address", func(value string, args []string) bool {
		_, err := net.ParseMAC(value)
		return err == nil
	})
	AddStringRule("not_regex", func(value string, args []string) bool {
		if len(args) < 1 {
			log.Print("not_regex must have 1 argument")
			return true
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			log.Printf("not_regex arg is not valid regex: %v", err)
			return true
		}
		return !re.MatchString(value)
	})
	AddStringRule("regex", func(value string, args []string) bool {
		if len(args) < 1 {
			log.Print("regex must have 1 argument")
			return true
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			log.Printf("regex arg is not valid regex: %v", err)
			return true
		}
		return re.MatchString(value)
	})
	AddStringRule("timezone", func(value string, args []string) bool {
		_, err := time.LoadLocation(value)
		return err == nil
	})
	AddStringRule("url", func(value string, args []string) bool {
		u, err := url.Parse(value)
		return err == nil && u.Host != "" && u.Scheme != ""
	})
	AddStringRule("uuid", func(value string, args []string) bool {
		_, err := uuid.Parse(value)
		return err == nil
	})
	AddStringRule("length", func(value string, args []string) bool {
		length, err := strconv.Atoi(args[0])
		if err != nil {
			log.Printf("length argument must be int, '%s' given", args[0])
			return true
		}
		return len(value) == length
	})
	AddStringRule("length_between", func(value string, args []string) bool {
		minLength, err := strconv.Atoi(args[0])
		if err != nil {
			log.Printf("length arguments must be int, '%s' given", args[0])
			return true
		}
		maxLength, err := strconv.Atoi(args[1])
		if err != nil {
			log.Printf("length arguments must be int, '%s' given", args[1])
			return true
		}
		return len(value) >= minLength && len(value) <= maxLength
	})
	AddStringRule("in", func(value string, args []string) bool {
		for _, arg := range args {
			if value == arg {
				return true
			}
		}
		return false
	})
	AddStringRule("not_in", func(value string, args []string) bool {
		for _, arg := range args {
			if value == arg {
				return false
			}
		}
		return true
	})
}

func AddStringRule(key string, cb func(value string, args []string) bool) {
	AddRule(key, func(options *ValidationOptions) bool {
		value, ok := options.Value.(string)
		if !ok {
			return true
		}
		return cb(value, options.Arguments)
	})
}

func alpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func numeric(c rune) bool {
	return c >= '0' && c <= '9'
}
