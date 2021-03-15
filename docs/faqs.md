## How do I share feedback on this plugin?

Feel free to create a GitHub issue or join the Bitbucket Plugin channel on our community Mattermost instance to discuss.

## How does the plugin save user data for each connected Bitbucket user?

Bitbucket user tokens are AES encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the PluginKeyValueStore table in your Mattermost database.
