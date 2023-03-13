package basic

import (
	"encoding/xml"
	"io/ioutil"
	"plugin"
	"strings"

	"github.com/oswaldoooo/octools/toolsbox"
)

type confinfo struct {
	XMLName xml.Name   `xml:"gocache"`
	Plugins plugininfo `xml:"plugins"`
	Paths   pathinfo   `xml:"paths"`
}
type plugininfo struct {
	XMLName xml.Name    `xml:"plugins"`
	Plugin  []theplugin `xml:"plugin"`
}
type theplugin struct {
	Class    string `xml:"classname"`
	FileName string `xml:"filename"`
}

var pluginlist = []theplugin{}

func LoadPlugin(plugin_name string) (pluginer *plugin.Plugin, err error) {
	pluginer, err = toolsbox.ScanPluginByName(plugin_name, ROOTPATH+"/plugins/")
	return
}
func ReadConf(filepath string) (err error) {
	if strings.Contains(filepath, ".xml") {
		content, err := ioutil.ReadFile(filepath)
		if err == nil {
			conf := new(confinfo)
			err = xml.Unmarshal(content, conf)
			if err == nil {
				for _, ve := range conf.Plugins.Plugin {
					switch strings.ToLower(ve.Class) {
					//支持的类
					case "fuzzy match":
						pluginlist = append(pluginlist, ve)
					}
				}
			}
		}
	}
	return
}

// load the pluginlist
func LoadPluginList() {
	var pluginer *plugin.Plugin
	var err error
	for _, vea := range pluginlist {
		switch vea.Class {
		case "fuzzy match":
			pluginer, err = LoadPlugin(vea.FileName)
			if err == nil {
				resfun, err := scanfuzzymatchfunc(pluginer)
				if err == nil {
					Fuzzy_Match_Func_Pool[vea.FileName] = resfun
				} else {
					errorlog.Println(err)
				}
			}
		}
	}
	pluginlist = nil
}
func scanfuzzymatchfunc(pluginer *plugin.Plugin) (resfunc func(target, tocompare string) bool, err error) {
	fzm, err := pluginer.Lookup("FuzzyMatch")
	if err == nil {
		resfunc = fzm.(func(target, tocompare string) bool)
	}
	return
}
