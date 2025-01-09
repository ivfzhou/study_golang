package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

func main() {}

func RunWithDefaultPort() {
	engine := gin.Default()

	// 默认使用 :8080 端口
	// t.Error(engine.Run())

	// 环境变量控制端口
	err := os.Setenv("PORT", "8081")
	if err != nil {
		log.Fatal(err)
	}

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	log.Fatalln(engine.Run())
}

func PathParam() {
	engine := gin.Default()

	// 匹配 /user/ivfzhou，不匹配 /user、/user/
	engine.GET("/user/:name", func(c *gin.Context) {
		// 获取路由参数
		log.Println("获取到的路由参数:name：", c.Param("name"))
		c.String(http.StatusOK, "命中路由 %s", c.FullPath())
	})

	// 匹配 /user/ivfzhou/、/user/ivfzhou/send
	// 如果没其他路由匹配 /user/ivfzhou 将重定向到 /user/ivfzhou/
	engine.GET("/user/:name/*action", func(c *gin.Context) {
		// 获取路由参数
		log.Println("获取到的路由参数*action：", c.Param("action"))

		// 获取路由原型
		log.Println(c.FullPath())

		c.String(http.StatusOK, "命中路由 %s", c.FullPath())
	})

	// 不管定义顺序先后，更明确的路由被优先匹配
	engine.GET("/user/ivfzhou", func(c *gin.Context) {
		c.String(http.StatusOK, "命中路由 %s", c.FullPath())
	})

	log.Fatalln(engine.Run())
}

func QueryParam() {
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		// 获取查询参数
		log.Println(c.Query("name"))
		// 当没有 name 查询参数字段时会返回 ivfzhou
		log.Println(c.DefaultQuery("name", "ivfzhou"))
		// 获取查询参数中的映射
		log.Println(c.QueryMap("name"))
		// 获取查询参数中的列表
		log.Println(c.QueryArray("name"))
	})

	log.Fatalln(engine.Run())
}

func MultipartAndUrlencodedParam() {
	engine := gin.Default()

	engine.POST("/", func(c *gin.Context) {
		// 获取 MultipartAndUrlencoded 参数
		log.Println(c.PostForm("name"))
		// 当没有 name 参数字段时会返回 ivfzhou
		log.Println(c.DefaultPostForm("name", "ivfzhou"))
		// 获取映射数据
		log.Println(c.PostFormMap("name"))
		// 获取列表数据
		log.Println(c.PostFormArray("name"))
	})

	log.Fatalln(engine.Run())
}

func MultipartFile() {
	engine := gin.Default()

	// 设置请求体最大大小，默认 32MiB
	engine.MaxMultipartMemory = 1

	engine.POST("/", func(c *gin.Context) {
		// 获取单一文件
		file, err := c.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}

		// 将文件保存到指定位置
		err = c.SaveUploadedFile(file, "tmp.txt")
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(file.Filename)
		log.Println(file.Header)
	})
	log.Fatalln(engine.Run())
}

func MultipartFiles() {
	engine := gin.Default()

	engine.POST("/", func(c *gin.Context) {
		// 获取多个文件
		form, err := c.MultipartForm()
		if err != nil {
			log.Println(err)
			return
		}
		files := form.File["file"]
		log.Println(len(files))
	})
	log.Fatalln(engine.Run())
}

func GroupRoute() {
	engine := gin.Default()

	{
		v1 := engine.Group("/v1")
		v1.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "命中路由：%s", c.FullPath())
		})
	}

	log.Fatalln(engine.Run())
}

func UseMiddleware() {
	// 包含了 Recovery 和 Logger 中间件
	// engine := gin.Default()

	engine := gin.New()

	// 捕获恐慌，然后响应 500 状态码
	engine.Use(gin.Recovery())

	// 日志写入 gin.DefaultWriter，gin.DefaultWriter 默认标准输出
	engine.Use(gin.Logger())

	// 自定义捕获恐慌
	engine.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		switch v := recovered.(type) {
		case string:
			c.String(http.StatusInternalServerError, v)
			c.Abort()
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	log.Fatalln(engine.Run())
}

func Log() {
	// 减少日志输出
	gin.SetMode(gin.ReleaseMode)

	// 设置日志输出
	/*f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)*/

	engine := gin.New()

	// 关闭日志颜色格式化
	// gin.DisableConsoleColor()
	// 强制使用日志着色
	// gin.ForceConsoleColor()
	// engine.Use(gin.Logger())

	// 自定义日志格式化
	/*engine.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s %s %s %v \"%s\" %d %d %s\n",
			params.TimeStamp.Format("2006/01/02 15:04:05"),
			params.Method,
			params.Path,
			params.ClientIP,
			params.Latency,
			params.Request.UserAgent(),
			params.StatusCode,
			params.BodySize,
			params.ErrorMessage,
		)
	}))*/

	// 设置跳过打印日志的路由匹配规则
	/*logConfig := gin.LoggerConfig{
		SkipPaths: []string{"/"},
		Skip: func(c *gin.Context) bool {
			if c.FullPath() == "/skip" {
				return true
			}
			return false
		},
	}
	engine.Use(gin.LoggerWithConfig(logConfig))*/

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "命中路由：%s", c.FullPath())
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "命中路由：%s", c.FullPath())
	})
	engine.GET("/skip", func(c *gin.Context) {
		c.String(http.StatusOK, "命中路由：%s", c.FullPath())
	})

	log.Fatalln(engine.Run())
}

type Address string

// 自定义处理接收参数
func (a *Address) UnmarshalParam(param string) error {
	log.Println(param)
	*a = Address(param)
	return nil
}

func Bind() {
	engine := gin.New()

	engine.POST("/:name", func(c *gin.Context) {
		type Req struct {
			// json => application/json
			// xml => application/xml
			// form => application/x-www-form-urlencoded 也绑定在查询参数中的
			// uri 路由参数
			// header 请求头字段
			Name string `header:"name" uri:"name" json:"name,omitempty" xml:"name" form:"name,default=ivfzhou" binding:"required"`
			// default 默认值，使用分号分隔。collection_format 分隔符，ssv 空格，tsv 制表符
			Hobbies []string `form:"hobbies,default=1 2" collection_format:"ssv"`
			// time_format 时间格式化接收，unixNano 纳秒时间戳，unix 秒时间戳，2006-01-02 时间格式化
			Birthday time.Time         `form:"birthday" time_format:"2006-01-02" time_utc:"8"`
			Address  Address           `form:"address"`
			Books    map[string]string `form:"books"`
			// 绑定文件
			File *multipart.FileHeader `form:"file"`
			// 可绑定内嵌 form 字段
			Nested struct {
				Age int `form:"age"`
			}
		}
		var req Req
		// 根据 Content-Type 确定如何解析请求参数。会读空 reqBody，再次调用触发 EOF。若是 GET 请求解析查询参数
		// JSON XML YAML TOML Header Query Uri
		// ShouldBindBodyWith 会将 reqBody 读入上下文中，可再次调用函数绑定参数
		err := c.ShouldBind(&req)
		if err != nil {
			var verr validator.ValidationErrors
			if errors.As(err, &verr) {
				if verr[0].StructField() == "Name" {
					c.String(http.StatusBadRequest, "名字不可为空")
					return
				}
			}
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		log.Println(req)
		c.JSON(http.StatusOK, req)
	})

	log.Fatalln(engine.Run())
}

type CustomBinder struct{}

func (*CustomBinder) Name() string {
	return "form"
}

func (*CustomBinder) Bind(req *http.Request, obj any) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	err = binding.MapFormWithTag(obj, req.Form, "url")
	if err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(obj)
}

func CustomBind() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		type Req struct {
			Name string `url:"name"`
		}
		var req Req
		err := c.ShouldBindWith(&req, &CustomBinder{})
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, &req)
	})
	log.Fatalln(engine.Run())
}

func Render() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		// XML JSON ProtoBuf TOML SecureJSON JSONP AsciiJSON PureJSON DataFromReader File FileFromFS
		// c.FileFromFS("main.go", http.Dir("./"))
		c.JSONP(http.StatusOK, "params")
	})
	log.Fatalln(engine.Run())
}

func ServeFile() {
	engine := gin.New()
	// engine.Static("/", "./")
	// engine.StaticFS("/", http.Dir("./"))
	engine.StaticFileFS("/gin/", "gin.http", http.Dir("./"))
	log.Fatalln(engine.Run())
}

func ServeHTML() {
	engine := gin.New()

	// 自定义分隔符
	// engine.Delims("{[{", "}]}")

	// engine.LoadHTMLGlob("./gin/templates/*")
	// engine.LoadHTMLFiles("./gin/templates/index.html")
	// engine.LoadHTMLGlob("./gin/templates/**/*")
	engine.SetHTMLTemplate(template.Must(template.ParseFiles("./gin/templates/index.html")))
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello",
		})
	})
	log.Fatalln(engine.Run())
}

func Redirect() {
	engine := gin.New()

	engine.POST("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://baidu.com")
	})
	engine.POST("/handle", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		engine.HandleContext(c)
	})

	log.Fatalln(engine.Run())
}

func BasicAuth() {
	engine := gin.New()

	engine.GET("/", gin.BasicAuth(gin.Accounts{
		"ivfzhou": "123456",
	}), func(c *gin.Context) {
		log.Println(c.MustGet(gin.AuthUserKey))
		c.String(http.StatusOK, "PASS")
	})

	log.Fatalln(engine.Run())
}

func CopyCtx() {
	engine := gin.New()

	engine.GET("/", func(c *gin.Context) {
		ctxCopy := c.Copy()
		c.Request.URL.Path = "/hello"
		log.Println(ctxCopy.Request.URL.Path)
		c.String(http.StatusOK, "PASS")
	})

	log.Fatalln(engine.Run())
}

func Listen() {
	engine := gin.New()

	// http.ListenAndServe(":8080", engine)

	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}
	log.Fatalln(server.ListenAndServe())
}

func RunTLS() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	log.Fatalln(autotls.Run(engine, "localhost"))
}

func Cert() {
	engine := gin.New()
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	log.Fatalln(autotls.RunWithManager(engine, &m))
}

//go:embed templates/*
var f embed.FS

func TmplInBin() {
	engine := gin.New()
	tmpls := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	engine.SetHTMLTemplate(tmpls)
	engine.GET("/", func(c *gin.Context) {
		// 读取文件
		// f.ReadFile("templates/index.html")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "hello",
		})
	})
	log.Fatalln(engine.Run())
}

func MultiServer() {
	var g errgroup.Group
	g.Go(func() error {
		err := (&http.Server{
			Addr:    ":8080",
			Handler: gin.New(),
		}).ListenAndServe()
		return err
	})
	g.Go(func() error {
		err := (&http.Server{
			Addr:    ":8081",
			Handler: gin.New(),
		}).ListenAndServe()
		return err
	})
	log.Fatalln(g.Wait())
}

func Shutdown() {
	engine := gin.New()
	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}
	go func() {
		log.Println(server.ListenAndServe())
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func DebugPrintRouteFunc() {
	engine := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	engine.GET("/", func(c *gin.Context) {

	})
	log.Fatalln(engine.Run())
}

func ClientIP() {
	engine := gin.New()

	// 设置信任的代理来源
	err := engine.SetTrustedProxies([]string{"192.168.137.1"})
	if err != nil {
		log.Fatalln(err)
	}

	// 设置取客户端 IP 的请求头字段，优先级比 SetTrustedProxies 高
	engine.TrustedPlatform = "X-Real_Ip"

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	log.Fatalln(engine.Run())
}

func Test() {
	engine := gin.New()

	// 在测试代码中
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	engine.ServeHTTP(w, r)
}
