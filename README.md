# oauth2-cli

*** support notice ***

I wrote this for me and my tool chain and offer very limited support. Use at your own risk.

---

Basically a wrapper around the golang oauth lib with cobra.

## Usage

`oauth2-cli token -c config.json`

That was my main intention to have. It saves the tokens in a session and tries to refresh if possible.


Check the config-examples dir for some possible configs.