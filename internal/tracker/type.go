package tracker

import "encoding/xml"

type Tracker struct {
	XMLName   xml.Name `xml:"trackerinfo"`
	Type      string   `xml:"type,attr"`
	ShortName string   `xml:"shortName,attr"`
	LongName  string   `xml:"longName,attr"`
	SiteName  string   `xml:"siteName,attr"`

	Settings  Settings  `xml:"settings"`
	Servers   Servers   `xml:"servers"`
	ParseInfo ParseInfo `xml:"parseinfo"`

	Config Config
}

type Settings struct {
	Description struct {
		Text string `xml:"text,attr"`
	} `xml:"description"`
	Passkey string `xml:"passkey"`
}

type Servers struct {
	Server []struct {
		Network    string `xml:"network,attr"`
		Names      string `xml:"serverNames,attr"`
		Channels   string `xml:"channelNames,attr"`
		Announcers string `xml:"announcerNames,attr"`
	} `xml:"server"`
}

type ParseInfo struct {
	LinePatterns struct {
		Extract struct {
			Regex struct {
				Value string `xml:"value,attr"`
			} `xml:"regex"`
			Vars struct {
				Var []struct {
					Name string `xml:"name,attr"`
				} `xml:"var"`
			} `xml:"vars"`
		} `xml:"extract"`
	} `xml:"linepatterns"`

	LineMatched struct {
		Var []struct {
			Name   string `xml:"name,attr"`
			String []struct {
				Value string `xml:"value,attr"`
			} `xml:"string"`
			Var []struct {
				Name string `xml:"name,attr"`
			} `xml:"var"`
		} `xml:"var"`
		If []struct {
			SrcVar string `xml:"srcvar,attr"`
			Regex  string `xml:"regex,attr"`
			Var    struct {
				Name   string `xml:"name,attr"`
				String []struct {
					Value string `xml:"value,attr"`
				} `xml:"string"`
				Var []struct {
					Name string `xml:"name,attr"`
				} `xml:"var"`
			} `xml:"var"`
		} `xml:"if"`
		Extract []struct {
			SrcVar   string `xml:"srcvar,attr"`
			Optional string `xml:"optional,attr"`
			Regex    struct {
				Value string `xml:"value,attr"`
			} `xml:"regex"`
			Vars struct {
				Var []struct {
					Name string `xml:"name,attr"`
				} `xml:"var"`
			} `xml:"vars"`
		} `xml:"extract"`
		VarReplace struct {
			Name    string `xml:"name,attr"`
			SrcVar  string `xml:"srcvar,attr"`
			Regex   string `xml:"regex,attr"`
			Replace string `xml:"replace,attr"`
		} `xml:"varreplace"`
		SetRegex struct {
			SrcVar   string `xml:"srcvar,attr"`
			Regex    string `xml:"regex,attr"`
			VarName  string `xml:"varName,attr"`
			NewValue string `xml:"newValue,attr"`
		} `xml:"setregex"`
		ExtractTags struct {
			SrcVar   string `xml:"srcvar,attr"`
			Split    string `xml:"split,attr"`
			SetVarIf []struct {
				VarName string `xml:"varName,attr"`
				Regex   string `xml:"regex,attr"`
			} `xml:"setvarif"`
			Regex struct {
				Value string `xml:"value,attr"`
			} `xml:"regex"`
		} `xml:"extracttags"`
	} `xml:"linematched"`
}
