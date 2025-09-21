package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsDirExist 判断文件夹是否存在
// 参数：
//   - folderPath 文件夹路径
//
// 返回值：
//   - bool 是否存在, true: 存在, false: 不存在
func IsDirExist(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsFileExist 判断文件是否存在
// 参数：
//   - filePath 文件路径
//
// 返回值：
//   - bool 是否存在, true: 存在, false: 不存在
func IsFileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetFileAbsPath 获取文件绝对路径
// 参数：
//   - inPath 输入路径
//   - filename 文件名
//
// 返回值：
//   - outPath 输出路径
//   - err 错误信息
func GetFileAbsPath(inPath string, filename ...string) (outPath string, err error) {
	if !filepath.IsAbs(inPath) {
		ex, err_ := os.Executable()
		if err_ != nil {
			err = fmt.Errorf("获取可执行文件路径失败: %w", err_)
		} else {
			inPath = filepath.ToSlash(filepath.Join(filepath.Dir(ex), inPath))
		}
	}

	if len(filename) == 0 {
		outPath = inPath
	} else {
		outPath = filepath.ToSlash(filepath.Join(inPath, filepath.Join(filename...)))
	}
	return
}

func RemoveAllData(path string, isCfg bool) {
	if isCfg {
		_ = os.RemoveAll(path)
		return
	}
	absPath, err := GetFileAbsPath(".", "local_dns_proxy.ini")
	if err == nil {
		_ = os.RemoveAll(absPath)
	}
}
