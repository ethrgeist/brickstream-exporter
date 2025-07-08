package models

import (
	"encoding/xml"
	"fmt"
	"time"
)

type MetricsV5 struct {
	XMLName    xml.Name   `xml:"Metrics"`
	SiteID     string     `xml:"SiteId,attr"`
	SiteName   string     `xml:"Sitename,attr"`
	DeviceID   string     `xml:"DeviceId,attr"`
	DeviceName string     `xml:"Devicename,attr"`
	DivisionID string     `xml:"DivisionId,attr"`
	Properties Properties `xml:"Properties"`
	ReportData ReportData `xml:"ReportData"`
}

type Properties struct {
	Version         int   `xml:"Version"`
	TransmitTime    int64 `xml:"TransmitTime"`
	TransmitTimeUTC time.Time
	MacAddress      string `xml:"MacAddress"`
	IPAddress       string `xml:"IpAddress"`
	HostName        string `xml:"HostName"`
	HTTPPort        int    `xml:"HttpPort"`
	HTTPSPort       int    `xml:"HttpsPort"`
	Timezone        int    `xml:"Timezone"`
	TimezoneName    string `xml:"TimezoneName"`
	DST             int    `xml:"DST"`
	TimezoneParsed  *time.Location
	HwPlatform      string `xml:"HwPlatform"`
	SerialNumber    string `xml:"SerialNumber"`
	DeviceType      int    `xml:"DeviceType"`
	SwRelease       string `xml:"SwRelease"`
}

type ReportData struct {
	Interval int      `xml:"Interval,attr"`
	Reports  []Report `xml:"Report"`
}

type Report struct {
	Date      string    `xml:"Date,attr"`
	DateLocal time.Time `xml:"-"`
	DateUTC   time.Time `xml:"-"`
	Objects   []Object  `xml:"Object"`
}

type Object struct {
	ID         string  `xml:"Id,attr"`
	DeviceID   string  `xml:"DeviceId,attr"`
	DeviceName string  `xml:"Devicename,attr"`
	ObjectType string  `xml:"ObjectType,attr"`
	Name       string  `xml:"Name,attr"`
	ExternalID string  `xml:"ExternalId,attr"`
	Counts     []Count `xml:"Count"`
}

type Count struct {
	StartTime           string `xml:"StartTime,attr"`
	StartTimeLocal      time.Time
	StartTimeUTC        time.Time
	EndTime             string `xml:"EndTime,attr"`
	EndTimeLocal        time.Time
	EndTimeUTC          time.Time
	UnixStartTime       int64 `xml:"UnixStartTime,attr"`
	UnixStartTimeParsed time.Time
	Enters              int `xml:"Enters,attr"`
	Exits               int `xml:"Exits,attr"`
	Status              int `xml:"Status,attr"`
}

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias Properties
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}

	p.TimezoneParsed = time.FixedZone(fmt.Sprintf("UTC%+d", p.Timezone+p.DST), (p.Timezone+p.DST)*3600)

	p.TransmitTimeUTC = time.Unix(p.TransmitTime, 0).In(p.TimezoneParsed)

	return nil
}

// in an ideal world this would implemented by UnmarshalXML bound to each of the structs
// but since we can not inject the timezone from properties into the report and counts
// this is the path of least resistance

func (m *MetricsV5) Process() {
	for ri, report := range m.ReportData.Reports {
		layoutDate := "2006-01-02"
		layoutTime := "15:04:05"

		dateLocal, err := time.ParseInLocation(layoutDate, report.Date, m.Properties.TimezoneParsed)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}

		m.ReportData.Reports[ri].DateLocal = dateLocal
		m.ReportData.Reports[ri].DateUTC = dateLocal.UTC()

		for _, object := range report.Objects {
			for ci, count := range object.Counts {
				// start
				startTimeLocal, err := time.ParseInLocation(layoutTime, count.StartTime, m.Properties.TimezoneParsed)
				if err != nil {
					fmt.Println("Error parsing start time:", err)
					return
				}

				startDateTimeLocal := time.Date(
					dateLocal.Year(),
					dateLocal.Month(),
					dateLocal.Day(),
					startTimeLocal.Hour(),
					startTimeLocal.Minute(),
					startTimeLocal.Second(),
					0,
					m.Properties.TimezoneParsed,
				)

				m.ReportData.Reports[ri].Objects[0].Counts[ci].StartTimeLocal = startDateTimeLocal
				m.ReportData.Reports[ri].Objects[0].Counts[ci].StartTimeUTC = startDateTimeLocal.UTC()

				// end
				endTimeLocal, err := time.Parse(layoutTime, count.EndTime)
				if err != nil {
					fmt.Println("Error parsing start time:", err)
					return
				}

				endDateTimeLocal := time.Date(
					dateLocal.Year(),
					dateLocal.Month(),
					dateLocal.Day(),
					endTimeLocal.Hour(),
					endTimeLocal.Minute(),
					endTimeLocal.Second(),
					0,
					m.Properties.TimezoneParsed,
				)

				m.ReportData.Reports[ri].Objects[0].Counts[ci].EndTimeLocal = endDateTimeLocal
				m.ReportData.Reports[ri].Objects[0].Counts[ci].EndTimeUTC = endDateTimeLocal.UTC()

				// unixstart
				t := time.Unix(count.UnixStartTime, 0).UTC()
				m.ReportData.Reports[ri].Objects[0].Counts[ci].UnixStartTimeParsed = t

			}
		}
	}
}
