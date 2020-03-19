package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"oauth2-cli/internal"
	"os"
	"strconv"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)


var (
	oauth2Config oauth2.Config
	configFile string
	config internal.Config
)

func init() {
	tokenCmd.Flags().StringVarP(&configFile, configFileFlag, "c", "", "oauth config to use")
	err := tokenCmd.MarkFlagRequired(configFileFlag)
	fatalOnErr(err)
	RootCmd.AddCommand(tokenCmd)
}


func initConfig() {
	var err error
	config, err = internal.ReadConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	oauth2Config = internal.ToOauth2Config(config)
}

func generateRandomString() string {
	n := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return strconv.Itoa(n)
}

var state = generateRandomString()

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "login with an oauth2-config",
	Long: "So, this is where the magic happens. You need to create a config file and add it with the config flag." +
		"It will fetch a (so called) token and save it as a session. If the token is already there it will read it from the session." +
		"If the token is refreshable it will try to refresh it (not yet implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		printToken(token())
	},
}

func printToken(token *oauth2.Token) {
	if token == nil {
		return
	}

	fmt.Print(token.AccessToken)
}

func token() *oauth2.Token {
	token, err := internal.LoadSession(config.Id)
	if err != nil {
		fatalOnErr(err)
	}

	if (token != nil) && !internal.SessionIsExpired(token) {
		onVerbose("using token from session")
		return token
	}

	if (token != nil) && internal.SessionIsExpired(token) && token.RefreshToken != "" {
		onVerbose("refreshing token from session")
		return refreshToken(token)
	}

	onVerbose("fetching new token")
	return handleTokenFetch()
}

func handleTokenFetch() *oauth2.Token {
	if config.Type == internal.Oauth2_grant_type_code {
		return login()

	}
	if config.Type == internal.Oauth2_grant_type_client_credentials {
		return clientCredentials()
	}
	log.Fatalln("unknown or missing grant_type")
	return nil
}

func clientCredentials() *oauth2.Token {
	ctx := context.Background()
	onVerbose("fetching token")
	clientCredentialsConfig := internal.ToOauth2ClientCredentialsConfig(config)
	token, err := clientCredentialsConfig.Token(ctx)
	fatalOnErr(err)
	onVerbose("got token: " + marshalStruct(token))
	err = internal.SaveSession(config.Id, token)
	fatalOnErr(err)
	return token
}

func login() *oauth2.Token {
	redirectUrl, err := url.Parse(oauth2Config.RedirectURL)
	fatalOnErr(err)
	authUrl := oauth2Config.AuthCodeURL(state)
	onVerbose("Visit the URL for the auth dialog: " + authUrl)
	err = browser.OpenURL(authUrl)
	if err != nil {
		log.Println(err)
	}
	onVerbose("registering callback on " + redirectUrl.EscapedPath())

	http.HandleFunc(redirectUrl.Path, callbackHandler)
	onVerbose("listening on " + redirectUrl.Host)
	log.Println(redirectUrl.Host)
	log.Fatal(http.ListenAndServe(redirectUrl.Host, nil))
	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)
	if errArray, ok := queryParts["error"]; ok {
		// this is currently not displayed. The whole server should be managed with a waitGroup
		// see https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
		errMsg := "<p>Something went wrong</p>" +
			"<p><strong>%s</strong></p>" +
			"<p>%s</p>"
		fmt.Fprintf(w, errMsg,errArray[0], queryParts["error_description"][0])
		fatalOnErr(errors.New(errArray[0]))
	}
	code := queryParts["code"][0]
	onVerbose("code: " + code)

	token, err := oauth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("state", state), oauth2.SetAuthURLParam("scope", oauth2Config.Scopes[0]))
	fatalOnErr(err)

	onVerbose("Token: " + token.AccessToken)
	err = internal.SaveSession(config.Id, token)
	if err != nil {
		fmt.Fprintf(w, "<p>something went wrong, check the logs</p>")
		fatalOnErr(err)
	}

	printToken(token)

	msg := "<p><strong>Success!</strong></p>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintf(w, msg)

	go os.Exit(0)
}

func refreshToken(token *oauth2.Token) *oauth2.Token{
	ctx := context.Background()
	newToken, err := oauth2Config.TokenSource(ctx, token).Token()
	if err != nil {
		fatalOnErr(err)
	}
	onVerbose("successfully refreshed token, saving...")
	err = internal.SaveSession(config.Id, newToken)
	fatalOnErr(err)

	return newToken
}