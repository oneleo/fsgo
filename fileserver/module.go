package fileserver

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Module contain http file server module.
type Module struct {
	Path string
}

// Compress module will zip the specified path.
func (s *Module) Compress(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("op") == "archive" {
		s.Path = r.FormValue("execute")
		zipPath := filepath.Clean(s.Path)

		var zipName string
		if zipPath == "." {
			zipName = "ROOT.zip"
		} else {
			zipName = filepath.Base(zipPath) + ".zip"
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="`+zipName+`"`)

		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		// 收集錯誤用，因 filepath.Walk 收到 err 會停止遍歷。
		var zipErrors []error

		filepath.Walk(zipPath, func(relatePath string, info os.FileInfo, err error) error {
			subPath := strings.TrimPrefix(relatePath, zipPath)
			// 根目錄為空，故要排除在壓縮檔之外
			if subPath == "" {
				return nil
			}
			hdr, err := zip.FileInfoHeader(info)
			if err != nil {
				zipErrors = append(zipErrors, err)
				return nil
			}
			// subPath 開頭為 / 或 \，要將之刪除，否則無法加入至壓縮檔。
			hdr.Name = strings.TrimLeft(subPath, `/\`)
			// 資料夾（非檔案）若最後方沒有加 / 或 \，會無法解壓縮此資料夾
			if info.IsDir() {
				hdr.Name = filepath.FromSlash(hdr.Name + "/")
			}

			if debug {
				fmt.Println("hdr.Name: ", hdr.Name)
			}

			// 設置壓縮方式
			hdr.Method = zip.Deflate
			// 在壓縮檔內建置對映空殼檔（尚無對映的資料）
			zw, err := zipWriter.CreateHeader(hdr)
			if err != nil {
				zipErrors = append(zipErrors, err)
				return nil
			}
			// 若此檔是捷徑，則取得它原來的檔名資訊
			if info.Mode()&os.ModeSymlink == os.ModeSymlink { // 是捷徑
				// 若為捷徑，則將其視為一般檔案，並將其指向的連結寫入此檔作為紀錄
				var symlink string
				symlink, err = os.Readlink(relatePath)
				if err != nil {
					zipErrors = append(zipErrors, err)
					return nil
				}
				zw.Write([]byte(symlink))
				// _, err = io.Copy(zw, ioutil.NopCloser(bytes.NewBuffer([]byte(symlink))))
				// if err != nil {
				// 	return err
				// }
			} else if !info.IsDir() { // 是檔案
				// rdc, err := os.Open(relatePath)
				// defer rdc.Close()
				// _, err = io.Copy(zw, rdc)
				// if err != nil {
				// 	return err
				// }
				// 將此檔的資料寫入對映空殼檔後壓縮
				data, err := ioutil.ReadFile(relatePath)
				if err != nil {
					zipErrors = append(zipErrors, err)
					return nil
				}
				zw.Write(data)
			} else if info.IsDir() { // 是資料夾
				// _, err = io.Copy(zw, ioutil.NopCloser(bytes.NewBuffer(nil)))
				// if err != nil {
				// 	return err
				// }
				// 資料夾無需寫入資料，即保留空殼檔即可
				//zw.Write([]byte(nil))
			}
			return nil
		})
		for _, e := range zipErrors {
			log.Println("zip err: ", e)
		}
		return
	}
	return
}
