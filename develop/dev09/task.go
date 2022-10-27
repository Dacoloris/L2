package dev09

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===
Реализовать утилиту wget с возможностью скачивать сайты целиком
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func saveSinglePage(path string) error {
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create("index.html")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

type Link struct {
	Path     string
	Dir      string
	Filename string
}

func newLink(u *url.URL) *Link {
	if u.String() == rootPath {
		return &Link{u.String(), rootDir, "index.html"}
	}

	if strings.HasSuffix(u.Path, "/") {
		l := Link{u.String(), u.Host + u.Path, u.RequestURI()[len(u.Path):]}
		if l.Filename == "" {
			l.Filename = "index.html"
		}
		return &l
	}

	s := strings.Split(u.Host+u.Path, "/")
	l := Link{u.String(), strings.TrimSuffix(u.Host+u.Path, s[len(s)-1]), s[len(s)-1] + u.RequestURI()[len(u.Path):]}
	l.Filename = strings.ReplaceAll(l.Filename, "/", "")
	if !strings.HasSuffix(l.Filename, ".html") {
		l.Filename = "'" + l.Filename + "'"
	}
	return &l
}

func (l *Link) getFilePath() string {
	dir := strings.TrimSuffix(l.Dir, "/")
	osDir, _ := os.Getwd()
	return "file://" + osDir + "/" + dir + "/" + l.Filename
}

func saveWebSite(path string) error {
	uLinks[path] = struct{}{}

	u, err := url.Parse(path)
	if err != nil {
		return err
	}

	rootPath = u.Scheme + "://" + u.Host
	rootDir = u.Host

	fmt.Println("downloading web-site:", rootPath)

	link := newLink(u)

	err = download(link, level)
	if err != nil {
		return err
	}

	return nil
}

func download(l *Link, lvl int) error {
	if lvl < 0 {
		return nil
	}

	fmt.Println("download:", l.Path)
	//fmt.Println("dir:", l.Dir)
	//fmt.Println("filename:", l.Filename)

	resp, err := http.Get(l.Path)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	body = bytes.ReplaceAll(body, []byte("&amp;"), []byte("&"))

	paths, err := getLinks(body, l.Path)
	if err != nil {
		return err
	}
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) > len(paths[j]) })

	pathPref := strings.TrimSuffix(l.Path, "/")

	var links []*Link
	for _, p := range paths {
		if strings.HasPrefix(p, "/") {
			p = pathPref + p
		}
		u, err := url.Parse(p)
		if err != nil {
			return err
		}
		links = append(links, newLink(u))
	}

	replaceLinks(&body, paths, links)

	err = savePage(body, l)
	if err != nil {
		return err
	}

	for _, link := range links {
		err := download(link, lvl-1)
		if err != nil {
			return fmt.Errorf("downloading %s failed: %s", link.Path, err)
		}
	}

	return nil
}

func getLinks(body []byte, path string) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var links []string
	f(doc, &links, path)

	return links, nil
}

func f(n *html.Node, links *[]string, path string) {
	path = strings.TrimSuffix(path, "/")

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && (strings.HasPrefix(a.Val, rootPath) || strings.HasPrefix(a.Val, "/")) {
				absPath := a.Val
				if strings.HasPrefix(absPath, "/") {
					absPath = path + absPath
				}
				if _, ok := uLinks[absPath]; !ok {
					uLinks[absPath] = struct{}{}
					*links = append(*links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f(c, links, path)
	}
}

func replaceLinks(body *[]byte, paths []string, links []*Link) {
	//s := string(*body)
	for i, l := range links {
		fmt.Println(l.Path)
		fmt.Println(l.Dir)
		fmt.Println(l.Filename)
		fmt.Println(l.getFilePath())
		//re := regexp.MustCompile(`\Q` + l.Path + `\E`)
		//fmt.Println(l.Path, "\n", re.Match(*body))
		//fmt.Println(l.Path, "\n", strings.Contains(s, l.Path))
		//s = strings.ReplaceAll(s, l.Path, l.getFilePath())
		*body = bytes.ReplaceAll(*body, []byte(paths[i]), []byte(l.getFilePath()))
	}
	//*body = []byte(s)
}

func savePage(body []byte, l *Link) error {
	err := os.MkdirAll(l.Dir, os.ModePerm)
	if err != nil {
		return err
	}

	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(l.Dir)
	if err != nil {
		return err
	}

	file, err := os.Create(l.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, bytes.NewReader(body))
	if err != nil {
		return err
	}

	return nil
}

var level int
var recursive bool

func init() {
	flag.IntVar(&level, "l", 1, "depth for recursive download")
	flag.BoolVar(&recursive, "r", false, "enable recursive download")
	flag.Parse()
}

var uLinks = make(map[string]struct{})
var rootPath string
var rootDir string

func main() {
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "no url provided")
		return
	}
	path := flag.Arg(0)

	if recursive {
		if level < 1 {
			fmt.Fprintf(os.Stderr, "level of recursion starts with 1, found: %d", level)
			return
		}

		if err := saveWebSite(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("downloading completed")
	} else {
		fmt.Println("downloading web-page:", path)
		if err := saveSinglePage(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("downloaded")
	}

}
