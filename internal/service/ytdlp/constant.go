package ytdlp

import (
	"path/filepath"

	"github.com/AmbitiousJun/live-server/internal/constant"
)

// 二进制文件存放根路径
var parentPath, _ = filepath.Abs(filepath.Join(constant.Dir_DataRoot, "yt-dlp"))
