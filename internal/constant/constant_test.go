package constant_test

import (
	"encoding/base64"
	"log"
	"os"
	"testing"
)

// TestTrasferHtmlTemplate 将 html 模板文件转换成 HTML 模板字符串
func TestTrasferHtmlTemplate(t *testing.T) {
	filePath := `/Users/ambitious/Desktop/code/go/live-server/static/config_page.html`

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("读取文件失败: %v", err)
	}

	log.Printf("转换结果: [%s]", base64.StdEncoding.EncodeToString(bytes))

}
