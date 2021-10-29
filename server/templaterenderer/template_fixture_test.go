package templaterenderer

import (
	"github.com/mattermost/mattermost-plugin-bitbucket/server/webhookpayload"
)

var mmUserBitbucketAccountID = "123"

var bitBucketAccountIDToUsernameMappingTestCallback BitBucketAccountIDToUsernameMappingCallbackType = func(accountID string) string {
	if accountID == mmUserBitbucketAccountID {
		return "testMmUser"
	}

	return ""
}

func getTestRepository() webhookpayload.Repository {
	repository := webhookpayload.Repository{FullName: "mattermost-plugin-bitbucket"}
	repository.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket"
	return repository
}

func getTestOwnerThatDoesntHaveAccount() webhookpayload.Owner {
	actor := webhookpayload.Owner{}
	actor.AccountID = "1230jdklf"
	actor.NickName = "testnickname"
	actor.Links.HTML.Href = "https://bitbucket.org/test-testnickname-url/"
	return actor
}

func getTestOwnerThatHasMmAccount() webhookpayload.Owner {
	owner := webhookpayload.Owner{}
	owner.AccountID = mmUserBitbucketAccountID
	owner.NickName = "mmUserBitbucketNickName"
	owner.Links.HTML.Href = "https://bitbucket.org/test-mmUserBitbucketNickName-url/"

	return owner
}

func getTestCommentWithMentionAboutMmUser() webhookpayload.Comment {
	comment := webhookpayload.Comment{}
	comment.Links.HTML.Href = "https://bitbucket.org/test-comment-link/"
	comment.Content.HTML = "<p>this issue should be fixed by <span class=\"ap-mention\"" +
		" data-atlassian-id=\"" + mmUserBitbucketAccountID + "\">@mmUserBitbucketNickname</span></p>"

	return comment
}

func getTestIssue() webhookpayload.Issue {
	issue := webhookpayload.Issue{}
	issue.ID = 1
	issue.Title = "README.md is outdated"
	issue.Content.HTML = "<p>README.md should be updated</p>"
	issue.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated"

	return issue
}

func getTestIssueCreatedPayload() webhookpayload.IssueCreatedPayload {
	return webhookpayload.IssueCreatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
	}
}

func getTestIssueUpdatedPayload() webhookpayload.IssueUpdatedPayload {
	return webhookpayload.IssueUpdatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
	}
}

func getTestIssueUpdatedPayloadWithStatusChange() webhookpayload.IssueUpdatedPayload {
	return webhookpayload.IssueUpdatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
		Changes: webhookpayload.IssueChanges{
			Status: webhookpayload.IssueChangesStatus{
				New: "closed",
			},
		},
	}
}

func getTestIssueCommentCreatedPayload() webhookpayload.IssueCommentCreatedPayload {
	return webhookpayload.IssueCommentCreatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
		Comment:    getTestCommentWithMentionAboutMmUser(),
	}
}

func getTestPullRequest() webhookpayload.PullRequest {
	pullRequest := webhookpayload.PullRequest{
		ID:    1,
		Title: "Test title",
	}
	pullRequest.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1"
	pullRequest.Rendered.Description.HTML = "<p>Test description <span class=\"ap-mention\"" +
		" data-atlassian-id=\"" + mmUserBitbucketAccountID + "\">@mmUserBitbucketNickname</span></p>"

	return pullRequest
}

func getTestPullRequestCreatedPayload() webhookpayload.PullRequestCreatedPayload {
	return webhookpayload.PullRequestCreatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestUpdatedPayload() webhookpayload.PullRequestUpdatedPayload {
	return webhookpayload.PullRequestUpdatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestApprovedPayload() webhookpayload.PullRequestApprovedPayload {
	return webhookpayload.PullRequestApprovedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestCommentCreatedPayload() webhookpayload.PullRequestCommentCreatedPayload {
	return webhookpayload.PullRequestCommentCreatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		Comment:     getTestCommentWithMentionAboutMmUser(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestDeclinedPayload() webhookpayload.PullRequestDeclinedPayload {
	return webhookpayload.PullRequestDeclinedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestMergedPayload() webhookpayload.PullRequestMergedPayload {
	return webhookpayload.PullRequestMergedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestUnapprovedPayload() webhookpayload.PullRequestUnapprovedPayload {
	return webhookpayload.PullRequestUnapprovedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestRepoPushChangeCommit1() webhookpayload.RepoPushChangeCommit {
	commit := webhookpayload.RepoPushChangeCommit{}
	commit.Author.User = getTestOwnerThatDoesntHaveAccount()
	commit.Hash = "dca5546b6b1419qf71adcada81b457cac3dcbdcd"
	commit.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd"
	commit.Message = "edit readme"

	return commit
}

func getTestRepoPushChangeCommit2() webhookpayload.RepoPushChangeCommit {
	commit := webhookpayload.RepoPushChangeCommit{}
	commit.Author.User = getTestOwnerThatHasMmAccount()
	commit.Hash = "fd84b9bfa6b415461f2c6559c262ed0826f76e08"
	commit.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/fd84b9bfa6b415461f2c6559c262ed0826f76e08"
	commit.Message = "edit something else"

	return commit
}

func getTestRepoPushChange() webhookpayload.RepoPushChange {
	pushChange := webhookpayload.RepoPushChange{}
	pushChange.Commits = append(pushChange.Commits, getTestRepoPushChangeCommit1())
	pushChange.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/master"
	pushChange.New.Name = "master"

	return pushChange
}

func getTestRepoPushPayloadWithOneCommit() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChange())
	pl.Repository.FullName = "mattermost/mattermost-plugin-bitbucket"

	return pl
}

func getTestRepoPushPayloadWithTwoCommits() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChange())
	pl.Push.Changes[0].Commits = append(pl.Push.Changes[0].Commits, getTestRepoPushChangeCommit2())
	pl.Repository.FullName = "mattermost/mattermost-plugin-bitbucket"

	return pl
}

func getTestRepoPushBranchCreated() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeBranchCreated())

	return pl
}

func getTestRepoPushTagCreated() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeTagCreated())

	return pl
}

func getTestRepoPushBranchDeleted() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeBranchDeleted())

	return pl
}

func getTestRepoPushTagDeleted() webhookpayload.RepoPushPayload {
	pl := webhookpayload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeTagDeleted())

	return pl
}

func getTestRepoPushChangeBranchCreated() webhookpayload.RepoPushChange {
	pushChange := webhookpayload.RepoPushChange{}
	pushChange.New.Type = "branch"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-new-branch"
	pushChange.New.Name = "test-new-branch"

	return pushChange
}

func getTestRepoPushChangeTagCreated() webhookpayload.RepoPushChange {
	pushChange := webhookpayload.RepoPushChange{}
	pushChange.New.Type = "tag"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-new-tag"
	pushChange.New.Name = "test-new-tag"

	return pushChange
}

func getTestRepoPushChangeBranchDeleted() webhookpayload.RepoPushChange {
	pushChange := webhookpayload.RepoPushChange{}
	pushChange.Old.Type = "branch"
	pushChange.Old.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-old-branch"
	pushChange.Old.Name = "test-old-branch"

	return pushChange
}

func getTestRepoPushChangeTagDeleted() webhookpayload.RepoPushChange {
	pushChange := webhookpayload.RepoPushChange{}
	pushChange.Old.Type = "tag"
	pushChange.Old.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-old-tag"
	pushChange.Old.Name = "test-old-tag"

	return pushChange
}
