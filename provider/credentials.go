package provider

import (
	"fmt"
	"net/url"
)

type CommonCredentials struct {
	URI      string `json:"uri"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type InfluxDBPrometheusBasicAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InfluxDBPrometheusRemoteCredentials struct {
	URL                  string                                 `json:"url"`
	BasicAuthCredentials InfluxDBPrometheusBasicAuthCredentials `json:"basic_auth"`
}

type InfluxDBPrometheusCredentials struct {
	RemoteRead  []InfluxDBPrometheusRemoteCredentials `json:"remote_read"`
	RemoteWrite []InfluxDBPrometheusRemoteCredentials `json:"remote_write"`
}

type InfluxDBCredentials struct {
	InfluxDBPrometheus *InfluxDBPrometheusCredentials `json:"prometheus,omitempty"`
}

type Credentials struct {
	CommonCredentials

	InfluxDBCredentials
}

func BuildCredentials(
	serviceType string,
	username string,
	password string,
	hostname string,
	port string,
) (Credentials, error) {
	credentials := Credentials{}

	credentials.URI = (&url.URL{
		Scheme: "https",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%s", hostname, port),
	}).String()

	credentials.Port = port
	credentials.Hostname = hostname
	credentials.Username = username
	credentials.Password = password

	if serviceType == "elasticsearch" {
		// nothing to do
	} else if serviceType == "influxdb" {
		addInfluxDBCredentials(&credentials)
	} else {
		return Credentials{}, fmt.Errorf("Unknown service type %s", serviceType)
	}

	return credentials, nil
}

func addInfluxDBCredentials(credentials *Credentials) {
	remoteReadURL := fmt.Sprintf(
		"https://%s:%s/api/v1/prom/read?db=defaultdb",
		credentials.Hostname, credentials.Port,
	)
	remoteWriteURL := fmt.Sprintf(
		"https://%s:%s/api/v1/prom/write?db=defaultdb",
		credentials.Hostname, credentials.Port,
	)

	credentials.InfluxDBPrometheus = &InfluxDBPrometheusCredentials{
		RemoteRead: []InfluxDBPrometheusRemoteCredentials{{
			URL: remoteReadURL,
			BasicAuthCredentials: InfluxDBPrometheusBasicAuthCredentials{
				Username: credentials.Username,
				Password: credentials.Password,
			},
		}},

		RemoteWrite: []InfluxDBPrometheusRemoteCredentials{{
			URL: remoteWriteURL,
			BasicAuthCredentials: InfluxDBPrometheusBasicAuthCredentials{
				Username: credentials.Username,
				Password: credentials.Password,
			},
		}},
	}
}
