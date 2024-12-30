<!-- markdownlint-disable MD013 MD033 MD041 -->

<img src="https://placecats.com/millie_neo/200/200" align="right" alt="Icon">

# Furry Garden Market Onboarding

Onboarding journey for customer acquisition when working from Farmers Markets!

## Install

First let's setup the development domain `market.furrygarden.io`:

```shell
echo "127.0.0.1  furrygarden.io" >> /etc/hosts
echo "127.0.0.1  market.furrygarden.io" >> /etc/hosts
brew install mkcert
brew install nss
mkcert -install
mkcert market.furrygarden.io
```

Now you have the domain setup in your Hosts file, and valid self-signed certificate.

```shell
make start
```

^ Starts the App with Air

Visit [https://market.furrygarden.io/login](https://market.furrygarden.io)

## Using the Tunnel

To allow Shopify to proxy the `/rest` API endpoints, run this command:

```shell
make start-tunnel
```

And make sure the _App Configuration_ uses that `serveo.net` URL here:
[https://partners.shopify.com/3017671/apps/200369307649/edit](https://partners.shopify.com/3017671/apps/200369307649/edit)

Update `Proxy URL` with the serveo URL like so:
`https://f316e6c145b889d21a521452dc045921.serveo.net`

## Crontab

### FG Onboarding App

Restart on reboot, and make sure every minute that it runs:

```bash
@reboot ~/apps/fgonboard_beta_app/restart.sh >> ~/logs/apps/fgonboard_beta_app/cron.log 2>&1
* * * * * pgrep fgonboarding > /dev/null || ~/apps/fgonboard_beta_app/restart.sh >> ~/logs/apps/fgonboard_beta_app/cron.log 2>&1
```
