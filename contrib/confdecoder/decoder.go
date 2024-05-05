package confdecoder

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"strings"
)

// NormalDecoder ... 修正了 kratos 在 Jetbrains IDE / Linux 环境下，基础配置文件热加载生成临时文件导致的报错问题
func NormalDecoder(src *config.KeyValue, target map[string]interface{}) error {
	if src.Format == "" {
		return nil
	}
	// 过滤文件结尾是临时文件的内容
	if strings.HasSuffix(src.Format, "~") {
		return nil
	}
	// 如果结尾是 .swp 则忽略
	if strings.HasSuffix(src.Format, "swp") {
		return nil
	}
	// src.Key is filename, you can also use regexp to filter it.
	if strings.HasPrefix(src.Key, ".") || strings.HasSuffix(src.Key, "~") {
		// do nothing, skiped file named src.Key
		return nil
	}

	codec := encoding.GetCodec(src.Format)
	if codec == nil {
		return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
	}
	return codec.Unmarshal(src.Value, &target)
}
