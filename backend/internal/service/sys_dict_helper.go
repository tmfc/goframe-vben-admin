package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

func resolveAcceptLanguage(ctx context.Context) string {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ""
	}
	header := strings.TrimSpace(req.Header.Get("Accept-Language"))
	if header == "" {
		return ""
	}
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}
	lang := strings.TrimSpace(parts[0])
	if lang == "" {
		return ""
	}
	if idx := strings.Index(lang, ";"); idx >= 0 {
		lang = strings.TrimSpace(lang[:idx])
	}
	return strings.ToLower(lang)
}

func parseLabelI18n(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var data map[string]string
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil
	}
	return data
}

func encodeLabelI18n(data map[string]string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func resolveDictLabel(label string, labelI18n map[string]string, lang string) string {
	if label == "" {
		return ""
	}
	if len(labelI18n) == 0 || lang == "" {
		return label
	}
	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		return label
	}
	if resolved, ok := lookupLabel(labelI18n, lang); ok {
		return resolved
	}
	if idx := strings.Index(lang, "-"); idx > 0 {
		if resolved, ok := lookupLabel(labelI18n, lang[:idx]); ok {
			return resolved
		}
	}
	return label
}

func lookupLabel(labelI18n map[string]string, lang string) (string, bool) {
	for key, value := range labelI18n {
		if strings.EqualFold(key, lang) {
			return value, true
		}
	}
	return "", false
}
