package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/happy-beer-api/internal/dispensers"
	"github.com/ldassonville/happy-beer-api/internal/records"
	"github.com/ldassonville/happy-beer-api/pkg/core/ginutils"
	"github.com/ldassonville/happy-beer-api/pkg/core/ioc"
	"github.com/sirupsen/logrus"
	logrusdd "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"k8s.io/utils/pointer"
	"log"
	"net/http"
	"os"
)

type APIConfigs struct {
	Name int  `yaml:"name,omitempty"`
	Port *int `yaml:"port,omitempty"`
}

type APISecrets struct {
}

func main() {
	injector := new(ioc.Injector)

	initApp(injector)

	start(injector)
}

func initApp(injector *ioc.Injector) {

	// initialize Datadog tracer
	tracer.Start()
	defer tracer.Stop()

	logrus.SetOutput(os.Stdout)
	logrus.AddHook(&logrusdd.DDContextLogHook{})

	// configuration & secrets
	apiConfigs := &APIConfigs{}
	//config.InitConfigs(apiConfigs)
	injector.Register(ioc.IocConfigs, apiConfigs)

	// register application secret
	apiSecrets := &APISecrets{}
	//config.InitSecrets(apiSecrets)
	injector.Register(ioc.IocSecrets, apiSecrets)

	// Register gin engine
	engine, err := ginutils.NewEngine()
	if err != nil {
		logrus.Fatal(err)
	}
	injector.Register(ioc.IocEngine, engine)

	records.Setup(injector)
	dispensers.Setup(injector)
}

func start(injector *ioc.Injector) {

	config := injector.Get(ioc.IocConfigs).(*APIConfigs)
	engine := injector.Get(ioc.IocEngine).(*gin.Engine)

	if config.Port == nil {
		config.Port = pointer.Int(9000)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", *config.Port),
		Handler: engine,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
