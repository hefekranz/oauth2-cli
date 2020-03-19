package internal

import (
    "encoding/json"
    "golang.org/x/oauth2"
    "io/ioutil"
    "log"
    "os"
    "time"
)

var SessionDir = "oauth2-cli"
var DefaultFileMode = os.FileMode(int(0770))

func GetSessionDir() string {
    dir, err  := os.UserConfigDir()
    if err != nil {
        log.Fatal(err)
    }
    return dir + "/" + SessionDir + "/"
}

func FileOrDirExists(filePath string) bool {
    if _, err := os.Stat(filePath); !os.IsNotExist(err) {
        return true
    }
    return false
}

func InitSessionDir() (err error) {
    return os.MkdirAll(GetSessionDir(), DefaultFileMode)
}

func LoadSession(id string) (token *oauth2.Token, err error) {
    fileName := GetSessionFilePath(id)
    if !FileOrDirExists(fileName) {
        return
    }
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        return
    }
    err = json.Unmarshal(data, &token)
    return
}

func SaveSession(id string, token *oauth2.Token) (err error) {
    fileName := GetSessionFilePath(id)
    data, err := json.Marshal(token)
    err = ioutil.WriteFile(fileName, data, DefaultFileMode)
    return
}

func GetSessionFilePath(id string) string {
    return GetSessionDir() + id
}

func SessionIsExpired(token *oauth2.Token) bool {
    return time.Now().After(token.Expiry)
}
