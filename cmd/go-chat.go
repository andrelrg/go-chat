package main

import (
	"github.com/getclasslabs/go-chat/internal"
	"github.com/getclasslabs/go-chat/internal/config"
	"github.com/getclasslabs/go-chat/internal/repositories"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConf "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

func main() {

	cfg := jaegerConf.Configuration{
		ServiceName: "go-chat",
		Sampler:     &jaegerConf.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    &jaegerConf.ReporterConfig{
			LogSpans: false,
		},
	}

	jLogger := jaegerLog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegerConf.Logger(jLogger),
		jaegerConf.Metrics(jMetricsFactory),
	)

	if err != nil {
		log.Fatal("failed to initialize tracer")
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()


	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config.Config)
	if err != nil {
		panic(err)
	}

	repositories.Start()
	config.Socket()

	s := internal.NewServer()
	log.Println("waiting routes...")
	log.Fatal(http.ListenAndServe(config.Config.Server.Port, s.Router))
}
