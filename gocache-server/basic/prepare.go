package basic

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/oswaldoooo/octools/toolsbox"
)

var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

type basicconf struct {
	XMLName   xml.Name   `xml:"gocache"`
	Plugins   plugininfo `xml:"plugins"`
	Paths     pathinfo   `xml:"paths"`
	Extension bool       `xml:"extensions"`
}
type pathinfo struct {
	XMLName     xml.Name      `xml:"paths"`
	Common_Path string        `xml:"common_path"`
	Other_Conf  otherconfinfo `xml:"other_conf"`
}
type otherconfinfo struct {
	XMLName    xml.Name `xml:"other_conf"`
	Conf_Paths []string `xml:"conf_path"`
}

func init() {
	fmt.Println("==========start init basic zone==========")
	content, err := ioutil.ReadFile(ROOTPATH + "/conf/conf.xml")
	if err != nil {
		fmt.Println("cant find the conf.xml")
		os.Exit(1)
	}
	bcinfo := new(basicconf)
	err = xml.Unmarshal(content, bcinfo)
	if err == nil {
		//读取配置文件，配置文件中可包含其他关于插件的配置文件
		IsOpenExtensions = bcinfo.Extension
		CommonPath = bcinfo.Paths.Common_Path
		if len(bcinfo.Paths.Other_Conf.Conf_Paths) > 0 {
			fmt.Println("read other conf files number", len(bcinfo.Paths.Other_Conf.Conf_Paths))
			read_other_conf(bcinfo.Paths.Other_Conf.Conf_Paths)
		}
	} else {
		fmt.Println("warn >> unmarshal basic zone failed", err)
	}
}

// 读取插件配置文件
func read_other_conf(filenames []string) {
	badfile := 0
	for _, filename := range filenames {
		if len(filename) > 5 && strings.Contains(filename, ".xml") {
			fullpath := ROOTPATH + "/conf/" + filename
			err := ReadConf(fullpath)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			badfile++
			if len(filename) > 0 {
				fmt.Println(filename, "is not correct xml file")
			}
		}
	}
	fmt.Println("bad configure file", badfile)
}
