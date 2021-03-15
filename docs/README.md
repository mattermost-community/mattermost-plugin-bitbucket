# Mattermost Bitbucket Plugin

A Bitbucket plugin for Mattermost. Based on the mattermost-plugin-bitbucket developed by jfrerich.

This guide is intended for Mattermost System Admins setting up the Bitbucket plugin and Mattermost users who want information about the plugin functionality. For more information about contributing to this plugin, visit the Development section.

The Mattermost Bitbucket plugin uses a webhook to connect your Bitbucket account to Mattermost to listen for incoming Bitbucket events. Events notifications are via DM in Mattermost. The Events don’t need separate configuration.

After your System Admin has configured the Bitbucket plugin, run `/bitbucket connect` in a Mattermost channel to connect your Mattermost and Bitbucket accounts.

Once connected, you'll have access to the following features:

* **Daily reminders:** The first time you log in to Mattermost each day, get a post letting you know what issues and pull requests need your attention.
* **Notifications:** Get a direct message in Mattermost when someone mentions you, requests your review, comments on, or modifies one of your pull requests/issues, or assigns you on Bitbucket.
* **Post actions:** Create a Bitbucket issue from a post or attach a post message to an issue. Hover over a post to reveal the post actions menu and click More Actions (...).
* **Sidebar buttons:** Stay up-to-date with how many reviews, assignments, and open pull requests you have with buttons in the Mattermost sidebar.
* **Slash commands:** Interact with the Bitbucket plugin using the /bitbucket slash command. Read more about slash commands here.
