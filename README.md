# Mattermost Bitbucket Plugin

[![Build Status](https://github.com/mattermost/mattermost-plugin-bitbucket/actions/workflows/ci.yml/badge.svg)](https://github.com/mattermost/mattermost-plugin-bitbucket/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/badge/Code%20Coverage-%-brightgreen?style=flat-square)](https://github.com/mattermost/mattermost-plugin-bitbucket)
[![Latest Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-bitbucket?style=flat-square)](https://github.com/mattermost/mattermost-plugin-bitbucket/releases)
[![Help Wanted](https://img.shields.io/github/issues/mattermost/mattermost-plugin-bitbucket?style=flat-square&color=brightgreen&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-bitbucket/issues?q=is%3Aissue+is%3Aopen+label%3A%22Help+Wanted%22)

## Help wanted tickets [[here]](https://github.com/mattermost/mattermost-plugin-bitbucket/issues?q=is%3Aissue+is%3Aopen+label%3A%22Help+Wanted%22)

## Contents

- [Overview](#overview)
- [Features](#features)
- [Admin Guide\*](docs\admin-guide.md)
- [End User Guide](#end-user-guide)
  - [Get Started](#get-started)
  - [Use the Plugin](#use-the-plugin)
  - [Slash Commands](#slash-commands)
  - [FAQ](#faq)
- [Contribute](#contribute)
  - [Development](#development)
  - [Test Your Changes](#test-your-changes)
- [Licence](#license)
- [Security Vulnerability Disclosure](#security-vulnerability-disclosure)
- [Get Help](#get-help)

## Overview

A Bitbucket plugin for Mattermost. Based on the [mattermost-plugin-bitbucket](https://github.com/jfrerich/mattermost-plugin-bitbucket) developed by [jfrerich](https://github.com/jfrerich).

## Features

The Bitbucket plugin features include:

- **Daily reminders:** The first time you log in to Mattermost each day, get a post letting you know what issues and pull requests need your attention.
- **Notifications:** Get a direct message in Mattermost when someone mentions you, requests your review, comments on, or modifies one of your pull requests/issues, or assigns you on Bitbucket.
- **Post actions:** Create a Bitbucket issue from a post or attach a post message to an issue. Hover over a post to reveal the post actions menu and select **More Actions \(...\)**.
- **Sidebar buttons:** Stay up-to-date with how many reviews, assignments, and open pull requests you have with buttons in the Mattermost sidebar.
- **Slash commands:** Interact with the Bitbucket plugin using the `/bitbucket` slash command.

![Bitbucket plugin screenshot](https://user-images.githubusercontent.com/45372453/97643091-114a1500-1a47-11eb-9863-2e0e308706ea.png)

## End User Guide

### Get Started

### Use the Plugin

### Slash commands

- **Subscribe to a respository:** Use `/bitbucket subscriptions add` to subscribe a Mattermost channel to receive notifications for new pull requests, issues, branch creation, and more in a Bitbucket repository.
  - For instance, to post notifications for issues, issue comments, and pull requests from mattermost/mattermost-server, use: `/bitbucket subscribe mattermost/mattermost-server issues,pulls,issue_comments`
- **Get to do items:** Use `/bitbucket todo` to get an ephemeral message with items to do in Bitbucket, including a list of assigned issues and pull requests awaiting your review.
- **Update settings:** Use `/bitbucket settings` to update your settings for notifications and daily reminders.

Run `/bitbucket help` to see what else the slash command can do.

### FAQ

#### How do I share feedback on this plugin?

Feel free to create a GitHub issue or join the Bitbucket Plugin channel on our community Mattermost instance to discuss.

#### How does the plugin save user data for each connected Bitbucket user?

Bitbucket user tokens are AES-encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.

## Contribute

### Development

This plugin contains both a server and web app portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

### Test Your Changes

## License

This repository is licensed under the [Apache 2.0 License](https://github.com/mattermost/mattermost-plugin-bitbucket/blob/master/LICENSE).

## Security Vulnerability Disclosure

## Get Help
