package main

import (
	"flag"
	"github.com/giantswarm/vault-exporter/log"
	"net/http"
	"os"

	"github.com/giantswarm/microerror"
	vault_api "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	listenAddress   = flag.String("web.listen-address", ":9410", "Address to listen on for web interface and telemetry.")
	metricsPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	vaultCACert     = flag.String("vault-tls-cacert", "", "The path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate.")
	vaultClientCert = flag.String("vault-tls-client-cert", "", "The path to the certificate for Vault communication.")
	vaultClientKey  = flag.String("vault-tls-client-key", "", "The path to the private key for Vault communication.")
	sslInsecure     = flag.Bool("insecure-ssl", false, "Set SSL to ignore certificate validation.")
)

const (
	namespace = "vault"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Vault successful.",
		nil, nil,
	)
	initialized = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "initialized"),
		"Is the Vault initialised (according to this node).",
		nil, nil,
	)
	sealed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "sealed"),
		"Is the Vault node sealed.",
		nil, nil,
	)
	standby = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "standby"),
		"Is this Vault node in standby.",
		nil, nil,
	)
	info = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "info"),
		"Version of this Vault node.",
		[]string{"version", "cluster_name", "cluster_id"}, nil,
	)
)

// Exporter collects Vault health from the given server and exports them using
// the Prometheus metrics package.
type Exporter struct {
	client *vault_api.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter() (*Exporter, error) {
	vaultConfig := vault_api.DefaultConfig()

	if *sslInsecure {
		tlsconfig := &vault_api.TLSConfig{
			Insecure: true,
		}
		err := vaultConfig.ConfigureTLS(tlsconfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	if *vaultCACert != "" || *vaultClientCert != "" || *vaultClientKey != "" {

		tlsconfig := &vault_api.TLSConfig{
			CACert:     *vaultCACert,
			ClientCert: *vaultClientCert,
			ClientKey:  *vaultClientKey,
			Insecure:   *sslInsecure,
		}
		err := vaultConfig.ConfigureTLS(tlsconfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	client, err := vault_api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		client: client,
	}, nil
}

// Describe describes all the metrics ever exported by the Vault exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- initialized
	ch <- sealed
	ch <- standby
	ch <- info
}

func bool2float(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// Collect fetches the stats from configured Vault and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	health, err := e.client.Sys().Health()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		log.Errorf("Failed to collect health from Vault server: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
	ch <- prometheus.MustNewConstMetric(
		initialized, prometheus.GaugeValue, bool2float(health.Initialized),
	)
	ch <- prometheus.MustNewConstMetric(
		sealed, prometheus.GaugeValue, bool2float(health.Sealed),
	)
	ch <- prometheus.MustNewConstMetric(
		standby, prometheus.GaugeValue, bool2float(health.Standby),
	)
	ch <- prometheus.MustNewConstMetric(
		info, prometheus.GaugeValue, 1, health.Version, health.ClusterName, health.ClusterID,
	)
}

func init() {
	prometheus.MustRegister(version.NewCollector("vault_exporter"))
}

func main() {
	err := mainE()
	if err != nil {
		panic(microerror.JSON(err))
	}
}

func mainE() error {
	if (len(os.Args) > 1) && (os.Args[1] == "version") {
		version.Print("vault_exporter")
		return nil
	}

	flag.Parse()

	log.Infof("Starting vault_exporter", version.Info())
	log.Infof("Build context", version.BuildContext())

	exporter, err := NewExporter()
	if err != nil {
		return microerror.Mask(err)
	}

	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
             <head><title>Vault Exporter</title></head>
             <body>
             <h1>Vault Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <h2>Build</h2>
             <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
             </body>
             </html>`))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Infof("Listening on", *listenAddress)

	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
