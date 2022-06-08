package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/messagebird/sachet"
	"github.com/messagebird/sachet/provider/aspsms"
	"github.com/messagebird/sachet/provider/cm"
	"github.com/messagebird/sachet/provider/esendex"
	"github.com/messagebird/sachet/provider/exotel"
	"github.com/messagebird/sachet/provider/freemobile"
	"github.com/messagebird/sachet/provider/ghasedak"
	"github.com/messagebird/sachet/provider/infobip"
	"github.com/messagebird/sachet/provider/kannel"
	"github.com/messagebird/sachet/provider/kavenegar"
	"github.com/messagebird/sachet/provider/mediaburst"
	"github.com/messagebird/sachet/provider/nowsms"
	"github.com/messagebird/sachet/provider/otc"
	"github.com/messagebird/sachet/provider/sap"
	"github.com/messagebird/sachet/provider/sfr"
	"github.com/messagebird/sachet/provider/sipgate"
	"github.com/messagebird/sachet/provider/smsc"
	"github.com/messagebird/sachet/provider/turbosms"
)

var (
	listenAddress = flag.String("listen-address", ":9876", "The address to listen on for HTTP requests.")
	configFile    = flag.String("config", "config.yaml", "The configuration file")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := LoadConfig(*configFile); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	app := handlers{}

	http.HandleFunc("/alert", app.Alert)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/-/reload", app.Reload)

	hc := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "sachet")

	http.HandleFunc("/-/live", hc.LiveEndpoint)
	http.HandleFunc("/-/ready", hc.ReadyEndpoint)

	if os.Getenv("PORT") != "" {
		*listenAddress = ":" + os.Getenv("PORT")
	}

	log.Printf("Listening on %s", *listenAddress)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

// receiverConfByReceiver loops the receiver conf list and returns the first instance with that name.
func receiverConfByReceiver(name string) *ReceiverConf {
	for i := range config.Receivers {
		rc := &config.Receivers[i]
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

func providerByName(name string) (sachet.Provider, error) {
	// TODO: use map of providers instead
	switch name {
	case "infobip":
		return infobip.NewInfobip(config.Providers.Infobip), nil
	case "kannel":
		return kannel.NewKannel(config.Providers.Kannel), nil
	case "kavenegar":
		return kavenegar.NewKaveNegar(config.Providers.KaveNegar), nil
	case "turbosms":
		return turbosms.NewTurbosms(config.Providers.Turbosms), nil
	case "smsc":
		return smsc.NewSmsc(config.Providers.Smsc), nil
	case "exotel":
		return exotel.NewExotel(config.Providers.Exotel), nil
	case "cm":
		return cm.NewCM(config.Providers.CM), nil

	case "otc":
		return otc.NewOTC(config.Providers.OTC), nil
	case "mediaburst":
		return mediaburst.NewMediaBurst(config.Providers.MediaBurst), nil
	case "freemobile":
		return freemobile.NewFreeMobile(config.Providers.FreeMobile), nil
	case "aspsms":
		return aspsms.NewAspSms(config.Providers.AspSms), nil
	case "sipgate":
		return sipgate.NewSipgate(config.Providers.Sipgate), nil
	case "nowsms":
		return nowsms.NewNowSms(config.Providers.NowSms), nil
	case "sap":
		return sap.NewSap(config.Providers.Sap), nil
	case "esendex":
		return esendex.NewEsendex(config.Providers.Esendex), nil
	case "ghasedak":
		return ghasedak.NewGhasedak(config.Providers.Ghasedak), nil
	case "sfr":
		return sfr.NewSfr(config.Providers.Sfr), nil
	}

	return nil, fmt.Errorf("%s: Unknown provider", name)
}

func errorHandler(w http.ResponseWriter, status int, err error, provider string) {
	w.WriteHeader(status)

	data := struct {
		Error   bool
		Status  int
		Message string
	}{
		true,
		status,
		err.Error(),
	}
	// respond json
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("marshalling error: " + err.Error())
	}

	if _, err := w.Write(body); err != nil {
		log.Fatalf("marshalling error: " + err.Error())
	}

	log.Println("error: " + string(body))
	requestTotal.WithLabelValues(strconv.FormatInt(int64(status), 10), provider).Inc()
}
