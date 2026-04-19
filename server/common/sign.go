package common

import (
	"strings"
	stdpath "path"

	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/setting"
	"github.com/alist-org/alist/v3/internal/sign"
)

func Sign(obj model.Obj, parent string, encrypt bool) string {
	if strings.HasSuffix(strings.ToLower(obj.GetName()), "cover.jpg") || strings.HasSuffix(strings.ToLower(obj.GetName()), "details.json") {
		return ""
	}
	if obj.IsDir() || (!encrypt && !setting.GetBool(conf.SignAll)) {
		return ""
	}
	return sign.Sign(stdpath.Join(parent, obj.GetName()))
}
