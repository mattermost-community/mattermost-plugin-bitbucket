package template_renderer

import "github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"

var mmUserBitbucketAccountId = "123"

var bitBucketAccountIDToUsernameMappingTestCallback BitBucketAccountIDToUsernameMappingCallbackType = func(accountID string) string {
	if accountID == mmUserBitbucketAccountId {
		return "testMmUser"
	}

	return ""
}

func getTestRepository() webhook_payload.Repository {
	repository := webhook_payload.Repository{FullName: "mattermost-plugin-bitbucket"}
	repository.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket"
	return repository
}

func getTestOwnerThatDoesntHaveAccount() webhook_payload.Owner {
	actor := webhook_payload.Owner{}
	actor.AccountId = "1230jdklf"
	actor.NickName = "testnickname"
	actor.Links.HTML.Href = "https://bitbucket.org/test-testnickname-url/"
	return actor
}

func getTestOwnerThatHasMmAccount() webhook_payload.Owner {
	owner := webhook_payload.Owner{}
	owner.AccountId = mmUserBitbucketAccountId
	owner.NickName = "mmUserBitbucketNickName"
	owner.Links.HTML.Href = "https://bitbucket.org/test-mmUserBitbucketNickName-url/"

	return owner
}

func getTestCommentWithMentionAboutMmUser() webhook_payload.Comment {
	comment := webhook_payload.Comment{}
	comment.Links.HTML.Href = "https://bitbucket.org/test-comment-link/"
	comment.Content.HTML = "<p>this issue should be fixed by <span class=\"ap-mention\"" +
		" data-atlassian-id=\"" + mmUserBitbucketAccountId + "\">@mmUserBitbucketNickname</span></p>"

	return comment
}

func getTestIssue() webhook_payload.Issue {
	issue := webhook_payload.Issue{}
	issue.ID = 1
	issue.Title = "README.md is outdated"
	issue.Content.HTML = "<p>README.md should be updated</p>"
	issue.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated"

	return issue
}

func getTestIssueCreatedPayload() webhook_payload.IssueCreatedPayload {
	return webhook_payload.IssueCreatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
	}
}

func getTestIssueUpdatedPayload() webhook_payload.IssueUpdatedPayload {
	return webhook_payload.IssueUpdatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
	}
}

func getTestIssueUpdatedPayloadWithStatusChange() webhook_payload.IssueUpdatedPayload {
	return webhook_payload.IssueUpdatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
		Changes: webhook_payload.IssueChanges{
			Status: webhook_payload.IssueChangesStatus{
				New: "closed",
			},
		},
	}
}

func getTestIssueCommentCreatedPayload() webhook_payload.IssueCommentCreatedPayload {
	return webhook_payload.IssueCommentCreatedPayload{
		Repository: getTestRepository(),
		Actor:      getTestOwnerThatHasMmAccount(),
		Issue:      getTestIssue(),
		Comment:    getTestCommentWithMentionAboutMmUser(),
	}
}

func getTestPullRequest() webhook_payload.PullRequest {
	pullRequest := webhook_payload.PullRequest{
		ID:    1,
		Title: "Test title",
	}
	pullRequest.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1"
	pullRequest.Rendered.Description.HTML = "<p>Test description <span class=\"ap-mention\"" +
		" data-atlassian-id=\"" + mmUserBitbucketAccountId + "\">@mmUserBitbucketNickname</span></p>"

	return pullRequest
}

func getTestPullRequestCreatedPayload() webhook_payload.PullRequestCreatedPayload {
	return webhook_payload.PullRequestCreatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestUpdatedPayload() webhook_payload.PullRequestUpdatedPayload {
	return webhook_payload.PullRequestUpdatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestApprovedPayload() webhook_payload.PullRequestApprovedPayload {
	return webhook_payload.PullRequestApprovedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestCommentCreatedPayload() webhook_payload.PullRequestCommentCreatedPayload {
	return webhook_payload.PullRequestCommentCreatedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		Comment:     getTestCommentWithMentionAboutMmUser(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestDeclinedPayload() webhook_payload.PullRequestDeclinedPayload {
	return webhook_payload.PullRequestDeclinedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestMergedPayload() webhook_payload.PullRequestMergedPayload {
	return webhook_payload.PullRequestMergedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestPullRequestUnapprovedPayload() webhook_payload.PullRequestUnapprovedPayload {
	return webhook_payload.PullRequestUnapprovedPayload{
		Actor:       getTestOwnerThatHasMmAccount(),
		PullRequest: getTestPullRequest(),
		Repository:  getTestRepository(),
	}
}

func getTestRepoPushChangeCommit1() webhook_payload.RepoPushChangeCommit {
	commit := webhook_payload.RepoPushChangeCommit{}
	commit.Author.User = getTestOwnerThatDoesntHaveAccount()
	commit.Hash = "dca5546b6b1419qf71adcada81b457cac3dcbdcd"
	commit.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd"
	commit.Message = "edit readme"

	return commit
}

func getTestRepoPushChangeCommit2() webhook_payload.RepoPushChangeCommit {
	commit := webhook_payload.RepoPushChangeCommit{}
	commit.Author.User = getTestOwnerThatHasMmAccount()
	commit.Hash = "fd84b9bfa6b415461f2c6559c262ed0826f76e08"
	commit.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/fd84b9bfa6b415461f2c6559c262ed0826f76e08"
	commit.Message = "edit something else"

	return commit
}

func getTestRepoPushChange() webhook_payload.RepoPushChange {
	pushChange := webhook_payload.RepoPushChange{}
	pushChange.Commits = append(pushChange.Commits, getTestRepoPushChangeCommit1())
	pushChange.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/master"
	pushChange.New.Name = "master"

	return pushChange
}

func getTestRepoPushPayloadWithOneCommit() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChange())
	pl.Repository.FullName = "mattermost/mattermost-plugin-bitbucket"

	return pl
}

func getTestRepoPushPayloadWithTwoCommits() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChange())
	pl.Push.Changes[0].Commits = append(pl.Push.Changes[0].Commits, getTestRepoPushChangeCommit2())
	pl.Repository.FullName = "mattermost/mattermost-plugin-bitbucket"

	return pl
}

func getTestRepoPushBranchCreated() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeBranchCreated())

	return pl
}

func getTestRepoPushTagCreated() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeTagCreated())

	return pl
}

func getTestRepoPushBranchDeleted() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeBranchDeleted())

	return pl
}

func getTestRepoPushTagDeleted() webhook_payload.RepoPushPayload {
	pl := webhook_payload.RepoPushPayload{}
	pl.Repository = getTestRepository()
	pl.Actor = getTestOwnerThatHasMmAccount()
	pl.Push.Changes = append(pl.Push.Changes, getTestRepoPushChangeTagDeleted())

	return pl
}

func getTestRepoPushChangeBranchCreated() webhook_payload.RepoPushChange {
	pushChange := webhook_payload.RepoPushChange{}
	pushChange.New.Type = "branch"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-new-branch"
	pushChange.New.Name = "test-new-branch"

	return pushChange
}

func getTestRepoPushChangeTagCreated() webhook_payload.RepoPushChange {
	pushChange := webhook_payload.RepoPushChange{}
	pushChange.New.Type = "tag"
	pushChange.New.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-new-tag"
	pushChange.New.Name = "test-new-tag"

	return pushChange
}

func getTestRepoPushChangeBranchDeleted() webhook_payload.RepoPushChange {
	pushChange := webhook_payload.RepoPushChange{}
	pushChange.Old.Type = "branch"
	pushChange.Old.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-old-branch"
	pushChange.Old.Name = "test-old-branch"

	return pushChange
}

func getTestRepoPushChangeTagDeleted() webhook_payload.RepoPushChange {
	pushChange := webhook_payload.RepoPushChange{}
	pushChange.Old.Type = "tag"
	pushChange.Old.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-old-tag"
	pushChange.Old.Name = "test-old-tag"

	return pushChange
}
