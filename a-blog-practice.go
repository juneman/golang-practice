package main

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

///////////////////////////////////////////////
func buildTopNav(section string, wr http.ResponseWriter) {
	fmt.Fprintf(wr, "<div style=\"margin: 0 auto; width: 600px;\">")
	fmt.Fprintf(wr, "<div style=\"margin-top: 20px;\">")

	actived := "<span style=\"background-color: #999; color: #f8f8f8;  margin-right:10px\">"
	unactive := "<span style=\"margin-right:10px;\">"
	if section == "home" {
		fmt.Fprintf(wr, actived+"       <a href=\"/\">Home</a></span>")
	} else {
		fmt.Fprintf(wr, unactive+"<a href=\"/\">Home</a></span>")
	}
	if section == "blog" {
		fmt.Fprintf(wr, actived+"<a href=\"/blog\">Blog</a></span>")
	} else {
		fmt.Fprintf(wr, unactive+"<a href=\"/blog\">Blog</a></span>")
	}
	if section == "about" {
		fmt.Fprintf(wr, actived+"<a href=\"/about\">About</a></span>")
	} else {
		fmt.Fprintf(wr, unactive+"  <a href=\"/about\">About</a></span>")
	}
	fmt.Fprintf(wr, "        </div>")
}

func buildBody(html string, wr http.ResponseWriter) {
	fmt.Fprintf(wr, html)
}
func buildFooter(wr http.ResponseWriter) {
	fmt.Fprintf(wr, "<br /><br /><div style=\"border-top: 1px solid #999; margin-top:20px;\">powered by Bin &copy;2014</div></div>")
}

////////////
type IndexController struct{}

func (this *IndexController) IndexAction(params []string, wr http.ResponseWriter, req *http.Request) {
	buildTopNav("home", wr)
	buildBody("<div><h3>Welcome to my blog.</h3></div>", wr)
	buildFooter(wr)
}

//---
type BlogController struct{}

func (this *BlogController) IndexAction(params []string, wr http.ResponseWriter, req *http.Request) {
	buildTopNav("blog", wr)

	body := "<div><ul>"
	body += "<li><a href=\"/blog/Detail/How to write hello world in Golang\">How to write hello world in Golang</a></li>"
	body += "<li><a href=\"/blog/Detail/Golang in China\">Golang in China</a></li>"
	body += "<li><a href=\"/blog/Detail/C and Golang\">C an Golang</a></li>"
	body += "</ul></div>"
	buildBody(body, wr)

	buildFooter(wr)
}

func (this *BlogController) DetailAction(params []string, wr http.ResponseWriter, req *http.Request) {
	buildTopNav("blog", wr)

	body := "<div>"
	body += "<h3>" + params[2] + "</h3>"
	body += "<p>About this topic <a href=\"http://golang.org\">" + params[2] + "</a>"
	body += "<p>Go is an open source programming language</p>"
	body += "<p>that makes it easy to build simple,</p><p>reliable, and efficient software.</p>"
	body += "<p><a href=\"http://golang.org\">click here to see more!!!</a></p>"
	body += "</div>"
	buildBody(body, wr)

	buildFooter(wr)
}

/////////////
type AboutController struct{}

func (this *AboutController) IndexAction(params []string, wr http.ResponseWriter, req *http.Request) {
	buildTopNav("about", wr)
	buildBody("<br /><div>This is a blog app build with <a href=\"http://golang.org\">Golang</a></div>", wr)
	buildFooter(wr)

}

/////////////
type ErrorController struct{}

func (this *ErrorController) IndexAction(params []string, wr http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wr, "<h1>Welcome to My Blog</h1>")
	fmt.Fprintf(wr, "<h3>404 error.</h3><br />")
	fmt.Fprintf(wr, "<div>")
	fmt.Fprintf(wr, "       <span><a href=\"/blog\">Blog</a></span>")
	fmt.Fprintf(wr, "         <span><a href=\"/about\">About</a></span>")
	fmt.Fprintf(wr, "        </div>")
}

/////////////
var (
	mupers = map[string]reflect.Value{
		"/":      reflect.ValueOf(new(IndexController)),
		"/blog":  reflect.ValueOf(new(BlogController)),
		"/about": reflect.ValueOf(new(AboutController)),
	}
)

//////////////////

func dispatch(params []string, params_len int, args ...interface{}) (result []reflect.Value) {

	in := make([]reflect.Value, len(args))
	for k, p := range args {
		in[k] = reflect.ValueOf(p)
	}

	cname := reflect.ValueOf(new(IndexController))
	if params_len > 0 {
		if c, ok := mupers["/"+params[0]]; ok {
			cname = c
		}
	}

	fname := "IndexAction"
	if params_len > 1 {
		fname = params[1] + "Action"
	}

	f := cname.MethodByName(fname)
	if (f == reflect.Value{}) {
		reflect.ValueOf(new(ErrorController)).MethodByName("indexAction")
		return
	}

	result = f.Call(in)
	return
}

////////////////

type Handler struct {
	*regexp.Regexp
}

func (h *Handler) getParams(p string) (params []string, c int) {
	match := h.FindStringSubmatch(p)
	params = make([]string, len(match))
	c = 0
	for i := 1; i < len(match); i++ {
		if match[i] == "" {
			break
		}
		if match[i] == "/" {
			continue
		}
		params[c] = match[i]
		c++
	}
	return
}

func (h *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	params, c := h.getParams(req.URL.Path)
	if c == 1 && params[0] == "favicon.ico" {
		return
	}
	dispatch(params, c, params, wr, req)
}

func main() {
	h := Handler{regexp.MustCompile("^/([a-zA-Z.]*)([/]*)([a-zA-Z.]*)([/]*)(.*)")}
	http.ListenAndServe(":80", &h)
}
