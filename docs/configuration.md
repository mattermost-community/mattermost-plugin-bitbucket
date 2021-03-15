## Configuration

Configuration is started in Bitbucket and completed in Mattermost.

### Step 1: Register an OAuth Application in Bitbucket
Go to https://bitbucket.org and log in.
Visit the Settings page for your organization.
Click the OAuth tab under Access Management.
Click the Add consumer button and set the following values:
Name: Mattermost Bitbucket Plugin - <your company name>.
Callback URL: https://your-mattermost-url.com/plugins/bitbucket/oauth/complete, replacing https://your-mattermost-url.com with your Mattermost URL.
URL: https://github.com/kosgrz/mattermost-plugin-bitbucket.
Set:
Account: Email and Read permissions.
Projects: Read permission.
Repositories: Read and Write permissions.
Pull requests: Read permission.
Issues: Read and write permissions.
Read and Write permissions on Issues and Pull requests and Read on Repositories and Account for this OAuth consumer account. 
Save.
the Key and Secret in the resulting screen.
Go to System Console > Plugins > Bitbucket and enter the Bitbucket OAuth Client ID and Bitbucket OAuth Client Secret you copied in a previous step.
Hit Save.

### Step 2: Create a Webhook in Bitbucket
You must create a webhook for each repository you want to receive notifications for or subscribe to.
Go to the Repository settings page of your Bitbucket organization you want to send notifications from, then select Webhooks in the sidebar.
Click Add Webhook.
Set the following values:
Title: Mattermost Bitbucket Webhook - <repository_name>, replacing repository_name with the name of your repository.
URL: https://your-mattermost-url.com/plugins/bitbucket/webhook, replacing https://your-mattermost-url.com with your Mattermost URL.
Select Choose from a full list of triggers.
Select:
Repository: Push.
Pull Request: Created, Updated, Approved, Approval removed, Merged, Declined, Comment created.
Issue: Created, Updated, Comment created.
Hit Save.
If you have multiple repositories, repeat the process to create a webhook for each repository.

### Step 3: Configure the Plugin in Mattermost
If you have an existing Mattermost user account with the name bitbucket, the plugin will post using the bitbucket account but without a BOT tag.
To prevent this, either:
Convert the bitbucket user to a bot account by running mattermost user convert bitbucket --bot in the CLI.
or
If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot 
account with the name bitbucket.
Generate a Key
Open System Console > Plugins > Bitbucket and do the following:
Generate a new value for At Rest Encryption Key.
(Optional) Bitbucket Organization: Lock the plugin to a single Bitbucket organization by setting this field to the name of your Bitbucket organization.
(Optional) Enable Private Repositories: Allow the plugin to receive notifications from private repositories by setting this value to true.
Hit Save.
Go to System Console > Plugins > Management and click Enable to enable the Bitbucket plugin.
You're all set!
