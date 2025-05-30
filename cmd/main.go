package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	monitoredServices = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pangolin_monitored_services_total",
		Help: "Total number of services with pangolin/enabled annotation set to true.",
	})
)

func init() {
	prometheus.MustRegister(monitoredServices)
}

func monitorServices() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	annotationKey := os.Getenv("ANNOTATION_KEY")
	annotationValue := os.Getenv("ANNOTATION_VALUE")

	for {
		services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error listing services: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		count := 0
		for _, service := range services.Items {
			annotations := service.ObjectMeta.Annotations
			if val, exists := annotations[annotationKey]; exists && val == annotationValue {
				count++
			}
		}

		monitoredServices.Set(float64(count))
		time.Sleep(30 * time.Second)
	}
}

func main() {
	go monitorServices()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
