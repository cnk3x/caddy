package k3xurl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	HeaderAccept         = "Accept"
	HeaderAcceptLanguage = "Accept-Language"
	HeaderAcceptEncoding = "Accept-Encoding"
	HeaderUserAgent      = "User-Agent"
	HeaderContentType    = "Content-Type"
	HeaderCacheControl   = "Cache-Control" // no-cache
	HeaderPragma         = "Pragma"        // no-cache
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

// With 以一些选项开始初始化请求器
func With(ctx context.Context, options ...Option) *Request {
	return (&Request{ctx: ctx}).With(options...)
}

// Default 默认的请求器，仅处理小于400的请求
func Default(ctx context.Context) *Request {
	return With(ctx).BeforeProcess(DefaultStatusCheck)
}

// 一些特定方法的定义
type (
	Option       = func(*Request) error                                   // 请求选项
	Body         = func() (contentType string, body io.Reader, err error) // 请求提交内容构造方法
	HeaderOption = func(headers http.Header)                              // 请求头处理
	Process      = func(resp *http.Response) error                        // 响应处理器
	ProcessMw    = func(next Process) Process                             // 响应预处理器
)

// Request 请求构造
type Request struct {
	ctx     context.Context        // Context
	options []func(*Request) error // options

	// request fields
	method    string         // 接口请求方法
	url       string         // 请求地址
	query     string         // 请求链接参数
	buildBody Body           // 请求内容
	headers   []HeaderOption // 请求头处理
	beforeMw  []ProcessMw    // 中间件

	// client fields
	tryTimes []time.Duration // 重试时间和时机
	client   *http.Client    // client
}

// With 增加选项
func (c *Request) With(options ...Option) *Request {
	c.options = append(c.options, options...)
	return c
}

// TryAt 失败重试，等待休眠时间
func (c *Request) TryAt(times ...time.Duration) *Request {
	c.tryTimes = times
	return c
}

/*设置请求*/

// Method 设置请求方法
func (c *Request) Method(method string) *Request {
	c.method = method
	return c
}

// Url 设置请求链接
func (c *Request) Url(url string) *Request {
	c.url = url
	return c
}

// Query 设置请求Query参数
func (c *Request) Query(query string) *Request {
	c.query = query
	return c
}

// Body 设置请求提交内容
func (c *Request) Body(body Body) *Request {
	c.buildBody = body
	return c
}

// Form 提交表单
func (c *Request) Form(form url.Values) *Request {
	return c.Body(func() (contentType string, body io.Reader, err error) {
		return "application/x-www-form-urlencoded; charset=utf-8", strings.NewReader(form.Encode()), nil
	})
}

// HeaderSet 设置请求头
func (c *Request) HeaderSet(k string, vs ...string) *Request {
	c.headers = append(c.headers, func(headers http.Header) { headers.Set(k, strings.Join(vs, ",")) })
	return c
}

// HeaderAdd 添加请求头
func (c *Request) HeaderAdd(k string, vs ...string) *Request {
	c.headers = append(c.headers, func(headers http.Header) {
		if k = http.CanonicalHeaderKey(k); k == HeaderAccept || k == HeaderAcceptEncoding {
			headers.Set(k, strings.Join(append(headers[http.CanonicalHeaderKey(k)], vs...), ","))
		} else {
			for _, v := range vs {
				headers.Add(k, v)
			}
		}
	})
	return c
}

// HeaderDel 删除请求头
func (c *Request) HeaderDel(keys ...string) *Request {
	c.headers = append(c.headers, func(headers http.Header) {
		for _, k := range keys {
			headers.Del(k)
		}
	})
	return c
}

/*处理响应*/

// BeforeProcess 在处理之前的预处理
func (c *Request) BeforeProcess(mws ...ProcessMw) *Request {
	c.beforeMw = append(c.beforeMw, mws...)
	return c
}

// Process 处理响应
func (c *Request) Process(process Process) error {
	if c.client == nil {
		c.client = &http.Client{}
	}

	for _, apply := range c.options {
		if err := apply(c); err != nil {
			return err
		}
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if c.method == "" {
		c.method = http.MethodGet
	}

	requestUrl := c.url
	if c.query != "" {
		if strings.Contains(requestUrl, "?") {
			requestUrl += "&" + c.query
		} else {
			requestUrl += "?" + c.query
		}
	}

	if c.buildBody == nil {
		c.buildBody = NoBody
	}

	var resp *http.Response
	for i := 0; i < len(c.tryTimes)+1; i++ {
		contentType, body, err := c.buildBody()
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(c.ctx, c.method, requestUrl, body)
		if err != nil {
			return err
		}

		if contentType != "" {
			req.Header.Set(HeaderContentType, contentType)
		}

		for _, headerOption := range c.headers {
			headerOption(req.Header)
		}

		if resp, err = c.client.Do(req); err != nil {
			var ne net.Error
			if i < len(c.tryTimes) && errors.As(err, &ne) {
				log.Printf("第%d次出错: %v, %s后重试", i+1, err, c.tryTimes[i])
				select {
				case <-c.ctx.Done():
					return err
				case <-time.After(c.tryTimes[i]):
					continue
				}
			}
			log.Printf("第%d次出错: %v, 返回错误", i+1, err)
			return err
		}
		break
	}

	respBody := resp.Body
	defer func(closer io.Closer) { _ = closer.Close() }(respBody)

	if process == nil {
		process = NoProcess
	}

	for _, before := range c.beforeMw {
		process = before(process)
	}
	return process(resp)
}

// ProcessBytes 处理响应字节
func (c *Request) ProcessBytes() (data []byte, err error) {
	err = c.Process(func(resp *http.Response) (ex error) {
		data, ex = io.ReadAll(resp.Body)
		return
	})
	return
}

// Download 下载到文件
func (c *Request) Download(fn string) (err error) {
	return c.Process(func(resp *http.Response) error {
		tempFn := fn + ".temp"
		if err := os.MkdirAll(filepath.Dir(tempFn), 0755); err != nil {
			return err
		}
		f, err := os.Create(tempFn)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, resp.Body)
		_ = f.Close()
		if err != nil {
			return err
		}
		return os.Rename(tempFn, fn)
	})
}

/* 设置客户端 */

// UseClient 使用的客户端定义
func (c *Request) UseClient(client *http.Client) *Request {
	c.client = client
	return c
}

// CookieEnabled 开关 Cookie
func CookieEnabled(enabled ...bool) Option {
	if len(enabled) == 0 || enabled[0] {
		jar, _ := cookiejar.New(nil)
		return Jar(jar)
	}
	return Jar(nil)
}

// Jar 设置Cookie容器
func Jar(jar http.CookieJar) Option {
	return func(c *Request) error {
		c.client.Jar = jar
		return nil
	}
}

// UseClient 使用自定义的HTTP客户端
func UseClient(client *http.Client) Option {
	return func(r *Request) error {
		r.client = client
		return nil
	}
}

// DefaultStatusCheck 默认状态检查，如果Status不小于400则直接报错
func DefaultStatusCheck(next Process) Process {
	return func(resp *http.Response) error {
		if resp.StatusCode >= 400 {
			return fmt.Errorf("%s", resp.Status)
		}
		return next(resp)
	}
}

// NoBody 空请求体
func NoBody() (contentType string, body io.Reader, err error) { return "", nil, nil }

// NoProcess 不处理
func NoProcess(*http.Response) error {
	return nil
}
