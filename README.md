# Mattermost Bitbucket Plugin

A Bitbucket plugin for Mattermost. Based on the [mattermost-plugin-bitbucket](https://github.com/jfrerich/mattermost-plugin-bitbucket) developed by [jfrerich](https://github.com/jfrerich).

## Feature summary

The Bitbucket plugin features include:

* **Daily reminders:** The first time you log in to Mattermost each day, get a post letting you know what issues and pull requests need your attention.
* **Notifications:** Get a direct message in Mattermost when someone mentions you, requests your review, comments on, or modifies one of your pull requests/issues, or assigns you on Bitbucket.
* **Post actions:** Create a Bitbucket issue from a post or attach a post message to an issue. Hover over a post to reveal the post actions menu and select **More Actions \(...\)**.
* **Sidebar buttons:** Stay up-to-date with how many reviews, assignments, and open pull requests you have with buttons in the Mattermost sidebar.
* **Slash commands:** Interact with the Bitbucket plugin using the `/bitbucket` slash command.

![Bitbucket plugin screenshot](https://user-images.githubusercontent.com/45372453/97643091-114a1500-1a47-11eb-9863-2e0e308706ea.png)

## Admin guide

This guide is intended for Mattermost System Admins setting up the Bitbucket plugin and Mattermost users who want information about the plugin functionality. For more information about contributing to this plugin, visit the Development section. The Mattermost Bitbucket plugin uses a webhook to connect your Bitbucket account to Mattermost to listen for incoming Bitbucket events. Events notifications are via Direct Message in Mattermost. The Events don’t need separate configuration.

### Prerequisites

This guide assumes:

* You have a Bitbucket account.
* You're a Mattermost System Admin.
* You're running Mattermost Server v5.25 or higher.

### Installation

#### Marketplace installation

1. Go to **Main Menu > Plugin Marketplace** in Mattermost.
2. Search for "Bitbucket" or find the plugin from the list.
3. Select **Install**.
4. When the plugin has downloaded and been installed, select **Configure**.

#### Manual installation

If your server doesn't have access to the internet, you can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-bitbucket/releases) and upload it to your server via **System Console > Plugin Management**. The releases on this page are the same used by the Marketplace. To learn more about how to upload a plugin, see [the documentation](https://developers.mattermost.com/integrate/plugins/using-and-managing-plugins/).

### Configuration

Configuration is started in Bitbucket and completed in Mattermost.

#### Step 1: Register an OAuth application in Bitbucket

1. Go to [https://bitbucket.org](https://bitbucket.org) and log in.
2. From your profile avatar in the bottom left, select the workspace in the **Recent workspaces** list or select **All workspaces** for a full list.
   * Select **Settings** in the left sidebar to open the workspace settings.
   * Under **Apps and features**, select **OAuth consumers**.
3. Select **Add consumer** and set the following values:
   * **Name:** `Mattermost Bitbucket Plugin - <your company name>`.
   * **Callback URL:** `https://your-mattermost-url.com/plugins/bitbucket/oauth/complete`, replacing `https://your-mattermost-url.com` with your Mattermost deployment's Site URL.
   * URL: [https://github.com/mattermost/mattermost-plugin-bitbucket](https://github.com/mattermost/mattermost-plugin-bitbucket).
4. Set:
   * **Account:** `Email` and `Read` permissions.
   * **Projects:** `Read` permission.
   * **Repositories:** `Read` and `Write` permissions.
   * **Pull requests:** `Read permission`.
   * **Issues:** `Read` and `Write` permissions.
5. Save the **Key** and **Secret** in the resulting screen.
6. Go to **System Console > Plugins > Bitbucket** 
7. Enter the Bitbucket **OAuth Client ID** and **Bitbucket OAuth Client Secret** you copied in a previous step.
8. Select **Save**.

#### Step 2: Create a webhook in Bitbucket

You must create a webhook for each repository you want to receive notifications for or subscribe to.

1. Go to the **Repository settings** page of the Bitbucket organization you want to send notifications from, then select **Webhooks** in the sidebar.
2. Select **Add Webhook**.
3. Set the following values:
   * **Title:** `Mattermost Bitbucket Webhook - <repository_name>`, replacing `repository_name` with the name of your repository.
   * **URL:** `https://your-mattermost-url.com/plugins/bitbucket/webhook?secret=SOME_SECRET`
      * replace `https://your-mattermost-url.com` with your Mattermost deployment's Site URL.
      * replace `SOME_SECRET` with the secret generated in System Console > Plugins > Bitbucket > Webhook Secret.
4. Select **Choose from a full list of triggers**.
5. Select:
   * **Repository:** `Push`.
   * **Pull Request:** `Created`, `Updated`, `Approved`, `Approval removed`, `Merged`, `Declined`, `Comment created`.
   * **Issue:** `Created`, `Updated`, `Comment created`.
6. Select **Save**.

If you have multiple repositories, repeat the process to create a webhook for each repository.

#### Step 3: Configure the plugin in Mattermost

If you have an existing Mattermost user account with the name `bitbucket`, the plugin will post using the `bitbucket` account but without a BOT tag.

To prevent this, either:

Convert the `bitbucket` user to a bot account by running `mattermost user convert bitbucket --bot` in the CLI.

or

If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot account with the name `bitbucket`.

#### Generate a key

Open **System Console > Plugins > Bitbucket** and do the following:

1. Generate a new value for **At Rest Encryption Key**.
2. \(Optional\) **Bitbucket Organization:** Lock the plugin to a single Bitbucket organization by setting this field to the name of your Bitbucket organization.
3. Select **Save**.
4. Go to **System Console > Plugins > Management** and select **Enable** to enable the Bitbucket plugin.

You're all set!

### Onboard users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their Bitbucket account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team, We've set up the Mattermost Bitbucket plugin, so you can get notifications in Mattermost. To get started, run the `/bitbucket connect` slash command from any channel within Mattermost to connect your Mattermost and Bitbucket accounts. Then, take a look at the slash commands section for details about how to use the plugin.

## User guide

### Slash commands

* **Subscribe to a respository:** Use `/bitbucket subscriptions add` to subscribe a Mattermost channel to receive notifications for new pull requests, issues, branch creation, and more in a Bitbucket repository.
  * For instance, to post notifications for issues, issue comments, and pull requests from mattermost/mattermost-server, use: `/bitbucket subscribe mattermost/mattermost-server issues,pulls,issue_comments`
* **Get to do items:** Use `/bitbucket todo` to get an ephemeral message with items to do in Bitbucket, including a list of assigned issues and pull requests awaiting your review.
* **Update settings:** Use `/bitbucket settings` to update your settings for notifications and daily reminders.

Run `/bitbucket help` to see what else the slash command can do.

### Frequently asked questions

#### How do I share feedback on this plugin?

Feel free to create a GitHub issue or join the Bitbucket Plugin channel on our community Mattermost instance to discuss.

#### How does the plugin save user data for each connected Bitbucket user?

Bitbucket user tokens are AES-encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.

## Development

This plugin contains both a server and web app portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

## License

This repository is licensed under the [Apache 2.0 License](https://github.com/mattermost/mattermost-plugin-bitbucket/blob/master/LICENSE).
