package main

import (
	goflag "flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// NvidiaCollector nvidia gpu collector
type NvidiaCollector struct {
	sync.Mutex
	numDevices  prometheus.Gauge
	memoryUsed  *prometheus.GaugeVec
	memoryTotal *prometheus.GaugeVec
	powerUsage  *prometheus.GaugeVec
	temperature *prometheus.GaugeVec
}

func newNvidiaCollector() *NvidiaCollector {
	namespace := "nvidia_gpu"
	labels := []string{"path", "uuid", "name"}
	nc := &NvidiaCollector{
		numDevices: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "num_devices",
				Help:      "Number of Nvidia GPU devices",
			},
		),
		memoryUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_used_mb",
				Help:      "Memory used by the GPU device in MB",
			},
			labels,
		),
		memoryTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_total_mb",
				Help:      "Memory Total of the GPU device in MB",
			},
			labels,
		),
		powerUsage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "power_usage",
				Help:      "Power usage of the GPU device in watts",
			},
			labels,
		),
		temperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "temperature",
				Help:      "Temperature of the GPU device in celsius",
			},
			labels,
		),
	}
	return nc
}

// Describe implement prometheus Collector interface
func (nc *NvidiaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.numDevices.Desc()
	nc.memoryUsed.Describe(ch)
	nc.memoryTotal.Describe(ch)
	nc.powerUsage.Describe(ch)
	nc.temperature.Describe(ch)
}

// Collect implement prometheus Collector interface
func (nc *NvidiaCollector) Collect(ch chan<- prometheus.Metric) {
	nc.Lock()
	defer nc.Unlock()

	nc.memoryUsed.Reset()
	nc.memoryTotal.Reset()
	nc.powerUsage.Reset()
	nc.temperature.Reset()

	count, err := nvml.GetDeviceCount()
	if err != nil {
		glog.Errorf("Failed to GetDeviceCount(): %v.", err)
		return
	}

	nc.numDevices.Set(float64(count))
	ch <- nc.numDevices

	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			glog.Warningf("Failed to NewDevice %d: %v.", i, err)
			continue
		}
		st, err := device.Status()
		if err != nil {
			glog.Warningf("Failed to get device[%d] status: %v.", i, err)
		}

		path := device.Path
		name := device.Model
		uuid := device.UUID

		nc.memoryUsed.WithLabelValues(path, uuid, *name).Set(float64(*st.Memory.Global.Used))
		nc.memoryTotal.WithLabelValues(path, uuid, *name).Set(float64(*device.Memory))
		nc.powerUsage.WithLabelValues(path, uuid, *name).Set(float64(*st.Power))
		nc.temperature.WithLabelValues(path, uuid, *name).Set(float64(*st.Temperature))

	}
	nc.memoryUsed.Collect(ch)
	nc.memoryTotal.Collect(ch)
	nc.powerUsage.Collect(ch)
	nc.temperature.Collect(ch)
}

// newCollectorCommand create and returns a new gpu collector
func newCollectorCommand() *cobra.Command {
	var (
		listenAddress string
		metricsPath   string
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run web interface",
		Long:  "Run web interface expose metrics",
		Run: func(cmd *cobra.Command, args []string) {
			goflag.Parse()
			glog.Infof("Starting gpu_exporter %s", version)

			if err := nvml.Init(); err != nil {
				glog.Errorf("Failed to initialize nvidia device: %v.", err)
				return
			}
			defer func() {
				glog.Info("Shudown NVML library.")
				nvml.Shutdown()
			}()

			prometheus.MustRegister(newNvidiaCollector())

			shudownC := make(chan struct{})
			go listenToSystemSignal(shudownC)

			http.Handle(metricsPath, promhttp.Handler())
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html>
						<head><title>GPU Exporter</title></head>
						<body>
						<h1>GPU Exporter</h1>
						<p><a href="` + metricsPath + `">Metrics</a></p>
						</body>
						</html>`))
			})

			go func() {
				http.ListenAndServe(listenAddress, nil)
			}()

			<-shudownC
			glog.Info("GPU exporter shutdown completed.")
			return

		},
	}

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.StringVar(&listenAddress, "web.listen-address", ":9470", "Address on which to expose metrics and web interface.")
	flags.StringVar(&metricsPath, "web.path", "/metrics", "Path under which to expose metrics.")

	return cmd
}

// listenToSystemSignal listen system signal and exit exporter.
func listenToSystemSignal(stopC chan<- struct{}) {
	glog.V(5).Info("Listen to system signal.")

	signalChan := make(chan os.Signal, 1)
	ignoreChan := make(chan os.Signal, 1)

	signal.Notify(ignoreChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		glog.V(3).Infof("GPU exporter shutdown by system signal: %s", sig)
		stopC <- struct{}{}
	}
}
