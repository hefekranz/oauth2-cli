package internal

import (
    "encoding/json"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/clientcredentials"
    "io/ioutil"
    "os"
)

const (
    Oauth2_grant_type_code = "authorization_code"
    Oauth2_grant_type_client_credentials = "client_credentials"
)

type Config struct {
    Id             string `json:"id,omitempty"`
    Type           string `json:"type,omitempty"`
    ClientId       string `json:"clientId,omitempty"`
    Oauth2AuthUrl  string `json:"oauth2AuthUrl,omitempty"`
    Oauth2TokenUrl string `json:"oauth2TokenUrl,omitempty"`
    RedirectUrl    string `json:"redirectUrl,omitempty"`
    Scope          string `json:"scope,omitempty"`
    ClientSecret   string `json:"clientSecret,omitempty"`
}

func ReadConfigFile(fileName string) (config Config, err error) {
    data, err := ReadFile(fileName)
    if err != nil {
        return
    }
    err = json.Unmarshal(data, &config)
    return
}

func ToOauth2Config(config Config) (c oauth2.Config) {
    c.ClientID = config.ClientId
    c.RedirectURL = config.RedirectUrl
    c.Endpoint.AuthURL = config.Oauth2AuthUrl
    c.Endpoint.TokenURL = config.Oauth2TokenUrl
    c.Scopes = []string{config.Scope}
    c.Endpoint.AuthStyle = 0
    c.ClientSecret = config.ClientSecret

    return
}

func ToOauth2ClientCredentialsConfig(config Config) (c clientcredentials.Config) {
    c.ClientID = config.ClientId
    c.TokenURL = config.Oauth2TokenUrl
    c.Scopes = []string{config.Scope}
    c.AuthStyle = 2
    c.ClientSecret = config.ClientSecret

    return
}

func ReadFile(fileName string) (data []byte, err error) {
    file, err := os.Open(fileName)
    if err != nil {
        return
    }
    return ioutil.ReadAll(file)
}
