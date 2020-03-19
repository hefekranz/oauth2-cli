package cmd

import (
    "encoding/json"
    "fmt"
    "github.com/spf13/cobra"
    "log"
    "oauth2-cli/internal"
    "os"
)

const name="oauth2-cli"

var (
    Version = "master"
    Date    = "undefined"
    Commit  = "undefined"
)

var Verbose bool

func init()  {
    RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
    err := internal.InitSessionDir()
    if err != nil {
        log.Fatal(err)
    }
}

var RootCmd = &cobra.Command{
    Use: name,
    Short: "fetches oauth2 access tokens for you and manages them as sessions.",
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}

func onVerbose(msg interface{})  {
    if Verbose {log.Println(msg)}
}

func fatalOnErr(err error)  {
    if err != nil {
        log.Fatal(err)
    }
}

func marshalStruct(thing interface{}) string {
    jsonThing, err := json.Marshal(thing)
    if err != nil {
        log.Println(err)
    }
    return string(jsonThing)
}


func printJson(data interface{}) {
    out, err := json.MarshalIndent(data,"","   ")
    if err != nil {
        fatalOnErr(err)
    }

    fmt.Print(string(out))
}