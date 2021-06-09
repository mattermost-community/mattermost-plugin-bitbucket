# Mattermost Bitbucket Plugin

A Bitbucket plugin for Mattermost. Based on the [mattermost-plugin-bitbucket](https://github.com/jfrerich/mattermost-plugin-bitbucket) developed by [jfrerich](https://github.com/jfrerich).

Visit the [Bitbucket documentation](https://mattermost.gitbook.io/bitbucket-plugin/) for guidance on installation, configuration, and usage.

## License

This repository is licensed under the [Apache 2.0 License](https://github.com/mattermost/mattermost-plugin-bitbucket/blob/master/LICENSE).

## About the Bitbucket Plugin

The Mattermost Bitbucket plugin uses a webhook to connect your Bitbucket account to Mattermost to listen for incoming Bitbucket events. Events notifications are via DM in Mattermost. The Events don’t need separate configuration. 

After your System Admin has [configured the Bitbucket plugin](#configuration), run `/bitbucket connect` in a Mattermost channel to connect your Mattermost and Bitbucket accounts.

Once connected, you'll have access to the following features:

* __Daily reminders__ - The first time you log in to Mattermost each day, get a post letting you know what issues and pull requests need your attention.
* __Notifications__ - Get a direct message in Mattermost when someone mentions you, requests your review, comments on or modifies one of your pull requests/issues, or assigns you on Bitbucket.
* __Post actions__ - Create a Bitbucket issue from a post or attach a post message to an issue. Hover over a post to reveal the post actions menu and click **More Actions (...)**.
* __Sidebar buttons__ - Stay up-to-date with how many reviews, assignments, and open pull requests you have with buttons in the Mattermost sidebar.
* __Slash commands__ - Interact with the Bitbucket plugin using the `/bitbucket` slash command. Read more about slash commands [here](#slash-commands).

## Before You Start

This guide assumes:

- You have a Bitbucket account.
- You're a Mattermost System Admin.
- You're running Mattermost v5.25 or higher.

## Configuration

Configuration is started in Bitbucket and completed in Mattermost. 

### Step 1: Register an OAuth Application in Bitbucket

1. Go to https://bitbucket.org and log in.
2. Visit the **Settings** page for your organization.
3. Click the **OAuth** tab under **Access Management**.
3. Click the **Add consumer** button and set the following values:
   - **Name:** `Mattermost Bitbucket Plugin - <your company name>`.
   - **Callback URL:** `https://your-mattermost-url.com/plugins/bitbucket/oauth/complete`, replacing `https://your-mattermost-url.com` with your Mattermost URL.
   - **URL:** `https://github.com/mattermost/mattermost-plugin-bitbucket`.
4. Set:
   - **Account:** `Email` and `Read` permissions.
   - **Projects:** `Read` permission.
   - **Repositories:** `Read` and `Write` permissions.
   - **Pull requests:** `Read` permission.
   - **Issues:** `Read` and `write` permissions.
5. Save the **Key** and **Secret** in the resulting screen.
6. Go to **System Console > Plugins > Bitbucket** and enter the **Bitbucket OAuth Client ID** and **Bitbucket OAuth Client Secret** you copied in a previous step.
7. Hit **Save**.

### Step 2: Create a Webhook in Bitbucket

You must create a webhook for each repository you want to receive notifications for or subscribe to.

1. Go to the **Repository settings** page of your Bitbucket organization you want to send notifications from, then select **Webhooks** in the sidebar.
2. Click **Add Webhook**.
3. Set the following values:
   - **Title:** `Mattermost Bitbucket Webhook - <repository_name>`, replacing `repository_name` with the name of your repository.
   - **URL:** `https://your-mattermost-url.com/plugins/bitbucket/webhook`, replacing `https://your-mattermost-url.com` with your Mattermost URL.
4. Select **Choose from a full list of triggers**.
5. Select:
   - **Repository:** `Push`.
   - **Pull Request:** `Created`, `Updated`, `Approved`, `Approval removed`, `Merged`, `Declined`, `Comment created`.
   - **Issue:** `Created`, `Updated`, `Comment created`.
6. Hit **Save**.

If you have multiple repositories, repeat the process to create a webhook for each repository.

### Step 3: Configure the Plugin in Mattermost

If you have an existing Mattermost user account with the name `bitbucket`, the plugin will post using the `bitbucket` account but without a `BOT` tag.

To prevent this, either:

- Convert the `bitbucket` user to a bot account by running `mattermost user convert bitbucket --bot` in the CLI.

or

- If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot 
account with the name `bitbucket`.

#### Generate a Key
  
Open **System Console > Plugins > Bitbucket** and do the following:

1. Generate a new value for **At Rest Encryption Key**.
2. (Optional) **Bitbucket Organization:** Lock the plugin to a single Bitbucket organization by setting this field to the name of your Bitbucket organization.
3. (Optional) **Enable Private Repositories:** Allow the plugin to receive notifications from private repositories by setting this value to `true`.
4. Hit **Save**.
5. Go to **System Console > Plugins > Management** and click **Enable** to enable the Bitbucket plugin.

You're all set!

## Using the Plugin

Once configuration is complete, run the `/bitbucket connect` slash command from any channel within Mattermost to connect your Mattermost account with Bitbucket.

## Onboarding Your Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their Bitbucket account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team, 

> We've set up the Mattermost Bitbucket plugin, so you can get notifications from Bitbucket in Mattermost. To get started, run the `/bitbucket connect` slash command from any channel within Mattermost to connect your Mattermost account with Bitbucket. Then, take a look at the [slash commands](#slash-commands) section for details about how to use the plugin.

## Slash Commands

* __Subscribe to a respository__ - Use `/bitbucket subscriptions add` to subscribe a Mattermost channel to receive notifications for new pull requests, issues, branch creation, and more in a Bitbucket repository.

   - For instance, to post notifications for issues, issue comments, and pull requests from `mattermost/mattermost-server`, use:
   ```
   /bitbucket subscribe mattermost/mattermost-server issues,pulls,issue_comments
   ```   
* __Get to do items__ - Use `/bitbucket todo` to get an ephemeral message with items to do in Bitbucket, including a list of assigned issues and pull requests awaiting your review.
* __Update settings__ - Use `/bitbucket settings` to update your settings for notifications and daily reminders.
* __And more!__ - Run `/bitbucket help` to see what else the slash command can do.

## Frequently Asked Questions

### How do I share feedback on this plugin?

Feel free to create a GitHub issue or [join the Bitbucket Plugin channel on our community Mattermost instance](https://community-release.mattermost.com/core/channels/plugin-bitbucket) to discuss.

### How does the plugin save user data for each connected Bitbucket user?

Bitbucket user tokens are AES encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.

## Development

This plugin contains both a server and web app portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.
