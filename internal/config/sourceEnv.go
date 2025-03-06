package config

import (
	"github.com/kyma-project/gpu-driver/internal/util"
	"github.com/tidwall/sjson"
	"regexp"
	"strings"
)

type sourceEnv struct {
	env          util.Environment
	fieldPath    string
	envVarPrefix string
}

func (s *sourceEnv) Read(in string) string {
	for k, v := range s.env.List() {
		if !strings.HasPrefix(k, s.envVarPrefix) {
			continue
		}

		fieldPath := s.fieldPath
		kk := strings.TrimPrefix(k, s.envVarPrefix)
		kk = strings.TrimPrefix(kk, "_")
		kk = strings.TrimPrefix(kk, "_")
		kk = strings.TrimPrefix(kk, "-")
		kk = strings.TrimPrefix(kk, "-")
		parts := strings.Split(strings.ToLower(kk), "__")
		rx := regexp.MustCompile("([-_][a-z])")
		for _, p := range parts {
			res := rx.ReplaceAllStringFunc(p, func(s string) string {
				s = strings.ToUpper(s)
				s = strings.ReplaceAll(s, "_", "")
				s = strings.ReplaceAll(s, "-", "")
				return s
			})
			fieldPath = ConcatFieldPath(fieldPath, res)
		}

		changed, err := sjson.Set(in, fieldPath, v)
		if err != nil {
			continue
		}
		in = changed
	}

	return in
}
