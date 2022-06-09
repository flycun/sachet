package main

import (
	"io/ioutil"

	"github.com/prometheus/alertmanager/template"
	"gopkg.in/yaml.v2"

	"github.com/flycun/sachet/provider/aspsms"
	"github.com/flycun/sachet/provider/cm"
	"github.com/flycun/sachet/provider/esendex"
	"github.com/flycun/sachet/provider/exotel"
	"github.com/flycun/sachet/provider/freemobile"
	"github.com/flycun/sachet/provider/ghasedak"
	"github.com/flycun/sachet/provider/infobip"
	"github.com/flycun/sachet/provider/kannel"
	"github.com/flycun/sachet/provider/kavenegar"
	"github.com/flycun/sachet/provider/mediaburst"
	"github.com/flycun/sachet/provider/nowsms"
	"github.com/flycun/sachet/provider/otc"
	"github.com/flycun/sachet/provider/sap"
	"github.com/flycun/sachet/provider/sfr"
	"github.com/flycun/sachet/provider/sipgate"
	"github.com/flycun/sachet/provider/smsc"
	"github.com/flycun/sachet/provider/turbosms"
)

type ReceiverConf struct {
	Name     string
	Provider string
	To       []string
	From     string
	Text     string
	Type     string
}

var config struct {
	Providers struct {
		Infobip      infobip.Config
		Kannel       kannel.Config
		KaveNegar    kavenegar.Config
		Exotel       exotel.Config
		CM           cm.Config
		Turbosms     turbosms.Config
		Smsc         smsc.Config
		OTC          otc.Config
		MediaBurst   mediaburst.Config
		FreeMobile   freemobile.Config
		AspSms       aspsms.Config
		Sipgate      sipgate.Config
		NowSms       nowsms.Config
		Sap          sap.Config
		Esendex      esendex.Config
		Ghasedak     ghasedak.Config
		Sfr          sfr.Config
	}

	Receivers []ReceiverConf
	Templates []string
}
var tmpl *template.Template

// LoadConfig loads the specified YAML configuration file.
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	tmpl, err = template.FromGlobs(config.Templates...)
	return err
}
