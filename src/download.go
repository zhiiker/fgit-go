package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func get(url, fpath string) {
	if url == "" || url == "--help" || url == "-h" {
		fmt.Println("" +
			"FastGit Get Command Line Tool\n" +
			"=============================\n" +
			"REMARKS\n" +
			"    Download with FastGit automatically\n" +
			"SYNTAX\n" +
			"    fgit [--help|-h]\n" +
			"    fgit get [URL<string>]\n" +
			"    fgit get [URL<string>] [Path<string>]\n" +
			"ALIASES\n" +
			"    fgit dl\n" +
			"    fgit download\n" +
			"EXAMPLE\n" +
			"    fgit get https://github.com/fastgitorg/fgit-go/archive/master.zip")
		os.Exit(0)
	} else {
		downloadFile(url, fpath)
	}
}

func downloadFile(url, fpath string) {
	urlSplit := strings.Split(url, "/")
	filename := urlSplit[len(urlSplit)-1]
	if fpath == "" {
		downloadFile(url, filename)
	}

	if isExists(fpath) {
		if isDir(fpath) {
			fpath = path.Join(fpath, filename)
			downloadFile(url, fpath)
		} else {
			isContinue := ' '
			fmt.Print("File with the same name exists. New file will cover the old file.\nDo you want to continue? [Y/n]")
			fmt.Scanf("%c", &isContinue)

			switch strings.ToLower(string(isContinue)) {
			case "y":
				os.Remove(fpath)
				goto startDown
			case "n":
				fmt.Println("User cancle the operation.")
				os.Exit(0)
			default:
				fmt.Println("Unknown input, exiting...")
				os.Exit(1)
			}

		}
	}

startDown:
	if strings.HasPrefix(url, "https://github.com/") {
		query := strings.Replace(url, "https://github.com/", "", -1)
		querySplit := strings.Split(query, "/")
		if len(querySplit) > 3 {
			// Source -> fastgitorg/fgit-go/blob/master/fgit.go
			// Target -> fastgitorg/fgit-go/master/fgit.go
			if querySplit[2] == "blob" {
				url = "https://raw.fastgit.org/"
				for _i, _s := range querySplit {
					if _i != 2 {
						// not /blob/
						if _i == len(querySplit)-1 {
							url += _s
						} else {
							url += _s + "/"
						}
					}
				}
				fmt.Println("Redirect ->", url)
			}
		}
	}
	newURL := strings.Replace(url, "https://github.com", "https://download.fastgit.org", -1)
	if newURL != url {
		fmt.Println("Redirect ->", newURL)
	}
	fmt.Println("File ->", fpath)
	fmt.Println("Downloading...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", newURL, nil)
	checkErr(err, "Http.Get create failed!", 1)
	req.Header.Set("User-Agent", "fgit/"+version)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	checkErr(err, "Http request failed!", 1)

	out, err := os.Create(fpath)
	checkErr(err, "File create failed!", 1)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	checkErr(err, "io.Copy failed!", 1)
	fmt.Println("Finished.")
	os.Exit(0)
}