// 代理是一种结构型设计模式，让你能提供真实服务对象的替代品给客户端使用。代理接收客户端的请求并进行一些处理（访问控制和缓存等），然后再将请求传递给服务对象。
// 延迟初始化（虚拟代理）。如果你有一个偶尔使用的重量级服务对象，一直保持该对象运行会消耗系统资源时，可使用代理模式。
// 访问控制（保护代理）。如果你只希望特定客户端使用服务对象，这里的对象可以是操作系统中非常重要的部分，而客户端则是各种已启动的程序（包括恶意程序），此时可使用代理模式。
// 本地执行远程服务（远程代理）。适用于服务对象位于远程服务器上的情形。
// 记录日志请求（日志记录代理）。适用于当你需要保存对于服务对象的请求历史记录时。
// 缓存请求结果（缓存代理）。适用于需要缓存客户请求结果并对缓存生命周期进行管理时，特别是当返回结果的体积非常大时。
// 智能引用。可在没有客户端使用某个重量级对象时立即销毁该对象。
package design_pattern_test

import (
	"fmt"
	"testing"
)

func TestProxy(t *testing.T) {
	nginxServer := NewNginxServer()
	appStatusURL := "/app/status"
	createUserURL := "/create/user"

	httpCode, body := nginxServer.HandleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.HandleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.HandleRequest(appStatusURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.HandleRequest(createUserURL, "POST")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)

	httpCode, body = nginxServer.HandleRequest(createUserURL, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", appStatusURL, httpCode, body)
}

// ===

type Server interface {
	HandleRequest(string, string) (int, string)
}

// =

type Nginx struct {
	application       *Application
	maxAllowedRequest int
	rateLimiter       map[string]int
}

func NewNginxServer() *Nginx {
	return &Nginx{
		application:       &Application{},
		maxAllowedRequest: 2,
		rateLimiter:       make(map[string]int),
	}
}

func (n *Nginx) HandleRequest(url, method string) (int, string) {
	allowed := n.checkRateLimiting(url)
	if !allowed {
		return 403, "Not Allowed"
	}
	return n.application.HandleRequest(url, method)
}

func (n *Nginx) checkRateLimiting(url string) bool {
	if n.rateLimiter[url] == 0 {
		n.rateLimiter[url] = 1
	}
	if n.rateLimiter[url] > n.maxAllowedRequest {
		return false
	}
	n.rateLimiter[url] = n.rateLimiter[url] + 1
	return true
}

// =

type Application struct{}

func (a *Application) HandleRequest(url, method string) (int, string) {
	if url == "/app/status" && method == "GET" {
		return 200, "Ok"
	}

	if url == "/create/user" && method == "POST" {
		return 201, "User Created"
	}
	return 404, "Not Ok"
}
