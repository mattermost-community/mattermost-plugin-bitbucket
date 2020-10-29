# Mattermost Bitbucket Plugin

A Bitbucket plugin for Mattermost. Based on the [mattermost-plugin-bitbucket](https://github.com/jfrerich/mattermost-plugin-bitbucket) developed by [jfrerich](https://github.com/jfrerich).

![image](https://user-images.githubusercontent.com/45372453/97643091-114a1500-1a47-11eb-9863-2e0e308706ea.png)

## Features

* __Daily reminders__ - the first time you log in to Mattermost each day, get a post letting you know what issues and pull requests need your attention.
* __Notifications__ - get a direct message in Mattermost when someone mentions you, requests your review, comments on or modifies one of your pull
  requests/issues, or assigns you on Bitbucket.
* __Sidebar buttons__ - stay up-to-date with how many reviews, assignments and open pull requests you have with buttons in the Mattermost sidebar.
* __Slash commands__ - interact with the Bitbucket plugin using the `/bitbucket` slash command
    * __Subscribe to a respository__ - Use `/bitbucket subscribe` to subscribe a
      Mattermost channel to receive posts for new pull requests and/or issues
      in a Bitbucket repository
    * __Get to do items__ - Use `/bitbucket todo` to get an ephemeral message with
      items to do in Bitbucket
    * __Update settings__ - Use `/bitbucket settings` to update your settings for the plugin
    * __And more!__ - Run `/bitbucket help` to see what else the slash command can do
* __Supports Bitbucket Enterprise__ - Works with SaaS and Enterprise versions
  of Bitbucket (Enterprise support added in version 0.6.0)

## Installation - Bitbucket

__Requires Mattermost 5.2 or higher. If you're running Mattermost 5.6+, it is strongly recommended to use plugin version 0.7.1+__

1. Install the plugin
    1. Download the latest version of the plugin from the Bitbucket releases page
    2. In Mattermost, go the System Console -> Plugins -> Management
    3. Upload the plugin
2. Register a Bitbucket OAuth app
    1. Go to https://bitbucket.org
    2. Click Avatar -> Bitbucket settings -> Settings -> Access Management (OAuth) -> Add consumer
        1. Use "Mattermost Bitbucket Plugin - &#060;your company name>" as the name
        2. Use "https://github.com/mattermost/mattermost-plugin-bitbucket" as the URL
        3. Use "https://your-mattermost-url.com/plugins/bitbucket/oauth/complete" as the callback URL, replacing `https://your-mattermost-url.com` with your Mattermost URL. 
        4. Set `Read` permissions on Issues, Repositories and Account for this OAuth consumer account 
        5. Save and copy the Key and Secret
    3. In Mattermost, go to System Console -> Plugins -> Bitbucket 
        1. Fill in the Client ID (Key) and Client Secret and save the settings
3. Create a Bitbucket webhook
    1. In Mattermost, go to the System Console -> Plugins -> Bitbucket -> Regenerate 
        1. Copy the "Webhook Secret"
    2. Go to the settings page of your Bitbucket repository and click on "Webhooks" in the sidebar
        1. Click "Add webhook"
        2. Use "Mattermost Bitbucket Webhook - &#060;repository_name>" as the title, where &#060;repository_name is the name of your repository 
        3. Use "https://your-mattermost-url.com/plugins/bitbucket/webhook" as the URL, replacing `https://your-mattermost-url.com` with your Mattermost URL 
        4. Select "Choose from a full list of trigger" as your trigger
        5. Select Issues: "Created", "Updated", and "Comment Created" 
    3. Save the webhook
    4. __Note for each organization you want to receive notifications for or subscribe to, you must create a webhook__
4. Generate an at rest encryption key
    1. Go to the System Console -> Plugins -> Bitbucket and click "Regenerate" under "At Rest Encryption Key"
    2. Save the settings
5. (Optional) Lock the plugin to a Bitbucket organization
    1. Go to System Console -> Plugins -> Bitbucket and set the Bitbucket
      Organization field to the name of your Bitbucket organization
6. (Optional) Enable private repositories
    1. Go to System Console -> Plugins -> Bitbucket and set Enable Private Repositories to true
    2. Note that if you do this after users have already connected their
      accounts to Bitbucket they will need to disconnect and reconnect their accounts to be able to use private repositories
7. (Not yet supported) Set your Enterprise URLs
    1. Go to System Console -> Plugins -> Bitbucket and set the Enterprise Base
      URL and Enterprise Upload URL fields to your Bitbucket Enterprise URLs, ex: `https://bitbucket.example.com`
    2. The Base and Upload URLs are often the same
8. Enable the plugin 
    1. Go to System Console -> Plugins -> Management and click "Enable" underneath the Bitbucket plugin
9. Test it out
    2. In Mattermost, run the slash command `/bitbucket connect`

## NOTES : bitbucket does not user a webhook secret

## Developing 

This plugin contains both a server and web app portion.

Run `ngrok` command to expose localhost to the internet 
* User Forwarding Address as &#060;your-mattermost-url.com>

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make check-style` to check the style.

Use `make deploy` to deploy the plugin to your local server. Before running `make deploy` you need to enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode). Edit your server configuration as follows:
                                                                                                      
```json
{
  "ServiceSettings": {
      ...
      "EnableLocalMode": true,
      "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
  }
}
```

Alternatively, you can authenticate with the server's API with credentials:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=sysadmin
export MM_ADMIN_PASSWORD=sysadmin
```

## Frequently Asked Questions

### How do I connect a repository instead of an organization?

Set up your Bitbucket webhook from the repository instead of the organization. Notifications and subscriptions will then be sent only for repositories you create webhooks for.

The reminder and `/bitbucket todo` will still search the whole organization, but only list items assigned to you.

## Feedback and Feature Requests

Feel free to create a Bitbucket issue or [join the Bitbucket Plugin channel on
our community Mattermost
instance](https://pre-release.mattermost.com/core/channels/plugin-bitbucket) to discuss.
