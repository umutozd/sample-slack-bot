# Sample Slack Bot
This repository contains a sample Slack bot implementation in Go, with the following features enabled:
- event handling
- interactivity
- OAuth flow

### Quick Install
Here's a link from by you can install the bot to your workspace, though it will not be functional as the server is not hosted anywhere:

<a href="https://slack.com/oauth/v2/authorize?client_id=4301412868374.5739063168660&scope=team:read,chat:write,app_mentions:read,im:history&user_scope="><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcSet="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

## Endpoints
### /slack/install
Handles redirects from the OAuth flow and obtains the authorization code. It exchanges it for an access and a refresh token and saves them for further use.

### /slack/events
Handles incoming Slack events. Currently handled events are:
- `app_home_opened`
- `message`

### /slack/interactive
Handles all interactivity for the bot such as button clicks, modal view updates etc.

## Development
First things first, you need to expose your locally-running server to the web. You can use it simply by installing [ngrok](https://ngrok.com/) and running:
```bash
ngrok http 8080
```
which will give you a public URL that forwards all requests to your local server.

Once set up, go to https://api.slack.com/apps and create new bot from scratch. In the bots configuration page:
1. Go to `App Home` and enable Home Tab and allow users to chat with the bot.
2. Go to `OAuth & Permissions` and:
 - Add the following redirect URL: `<your-ngrok-url>/slack/install`
 - Opt in for OAuth flow
 - Add following scopes:
   - `app_mentions:read`
   - `chat:write`
   - `im:history`
   - `team:read`
3. Go to `Event Subscriptions` and:
  - Add the following request URL: `<your-ngrok-url>/slack/events`, which will send a challenge to your locally-running server.
  - Subscribe to the following bot events:
    - `app_home_opened`
    - `app_mention`
    - `message.im`
4. Go to `Interactivity & Shortcuts` and after enabling it, add the following request URL: `<your-ngrok-url>/slack/interactive`

## Deployment
Go to `Manage Distribution` page of your bot's configuration page and click on the `Add to Slack` button, which will redirect you to a consent page where you are asked for permission by Slack and to which workspace you want to install the bot. Once done, you will be redirected to the installed app's `About` tab. That's it!
