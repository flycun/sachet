package main

import (
	"io/ioutil"

	"github.com/prometheus/alertmanager/template"
	"gopkg.in/yaml.v2"

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
