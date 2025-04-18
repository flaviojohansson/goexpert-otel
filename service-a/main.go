package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServicebResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type CepInputDTO struct {
	CEP string `json:"cep" binding:"required,min=8"`
}

type CepOutputDTO struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

var (
	serviceName  = os.Getenv("OTEL_SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func cepHandler(c *gin.Context) {

	var cepInputDTO CepInputDTO

	if err := c.ShouldBindJSON(&cepInputDTO); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	// Extract the trace context from the incoming request
	// and create a new span for the request
	// _, span := tracer.Start(c.Request.Context(), "cepHandler", oteltrace.WithAttributes(attribute.String("cep", cepInputDTO.CEP)))
	// defer span.End()

	url := fmt.Sprintf("http://service-b:8081/temperatura?cep=%s", cepInputDTO.CEP)

	// Create a client with tracing
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Create the request
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Make the HTTP request to service B
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	var servicebResponse ServicebResponse

	if err := json.NewDecoder(resp.Body).Decode(&servicebResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	cepOutputDTO := CepOutputDTO(servicebResponse)

	c.JSON(http.StatusOK, cepOutputDTO)
}

func main() {

	// Initialize OpenTelemetry
	ctx := context.Background()
	_, shutdown, err := InitTracer(ctx, serviceName, collectorURL)

	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %s, %v", collectorURL, err)
	}
	defer shutdown(ctx)

	// Set up propagator (should already be in your InitTracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	r.POST("/", cepHandler)
	r.Run(":8080")

}
