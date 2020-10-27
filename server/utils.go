package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mattermost/mattermost-server/v5/utils"
	"io"
	"strings"
)

func getBaseURL() string {
	return "https://api.bitbucket.org/2.0"
}

func getYourOrgReposSearchQuery(organizationName string) string {
	return getBaseURL() + "/repositories/" + organizationName + "?role=member"
}

func getYourAllReposSearchQuery() string {
	return getBaseURL() + "/repositories?role=member"
}

func getYourAssigneeIssuesSearchQuery(userAccountID, repoFullName string) string {
	return getBaseURL() + "/repositories/" + repoFullName + "/issues?q=" +
		utils.UrlEncode("assignee.account_id=\""+userAccountID+"\" AND state!=\"closed\"")
}

func getYourAssigneePRsSearchQuery(userAccountID, repoFullName string) string {
	return getBaseURL() + "/repositories/" + repoFullName + "/pullrequests?q=" +
		utils.UrlEncode("reviewers.account_id=\""+userAccountID+"\" AND state=\"open\"")
}

func getYourOpenPRsSearchQuery(userAccountID, repoFullName string) string {
	return getBaseURL() + "/repositories/" + repoFullName + "/pullrequests?q=" +
		utils.UrlEncode("author.account_id=\""+userAccountID+"\" AND state=\"open\"")
}

func getSearchIssuesQuery(repoFullName, searchTerm string) string {
	return getBaseURL() + "/repositories/" + repoFullName + "/issues?q=" +
		utils.UrlEncode("title ~ \""+searchTerm+"\"") + "&sort=-updated_on"
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], msg)
	finalMsg := base64.URLEncoding.EncodeToString(ciphertext)
	return finalMsg, nil
}

func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multiple of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}

func parseOwnerAndRepoAndReturnFullAlso(full, baseURL string) (string, string, string) {
	if baseURL == "" {
		baseURL = BitbucketBaseURL
	}
	full = strings.TrimSuffix(strings.TrimSpace(strings.Replace(full, baseURL, "", 1)), "/")
	splitStr := strings.Split(full, "/")

	if len(splitStr) == 1 {
		owner := splitStr[0]
		return owner, owner, ""
	} else if len(splitStr) != 2 {
		return "", "", ""
	}
	owner := splitStr[0]
	repo := splitStr[1]

	return fmt.Sprintf("%s/%s", owner, repo), owner, repo
}

func parseOwnerAndRepo(full, baseURL string) (string, string) {
	if baseURL == "" {
		baseURL = BitbucketBaseURL
	}
	full = strings.TrimSuffix(strings.TrimSpace(strings.Replace(full, baseURL, "", 1)), "/")
	splitStr := strings.Split(full, "/")

	if len(splitStr) == 1 {
		owner := splitStr[0]
		return owner, ""
	} else if len(splitStr) != 2 {
		return "", ""
	}
	owner := splitStr[0]
	repo := splitStr[1]

	return owner, repo
}

// getToDoDisplayText returns the text to be displayed in todo listings.
func getToDoDisplayText(baseURL, title, url, notifType string) string {
	owner, repo := parseOwnerAndRepo(url, baseURL)
	repoURL := fmt.Sprintf("%s%s/%s", baseURL, owner, repo)
	repoWords := strings.Split(repo, "-")
	if len(repo) > 20 && len(repoWords) > 1 {
		repo = "..." + repoWords[len(repoWords)-1]
	}
	repoPart := fmt.Sprintf("[%s/%s](%s)", owner, repo, repoURL)

	if len(title) > 80 {
		title = strings.TrimSpace(title[:80]) + "..."
	}
	titlePart := fmt.Sprintf("[%s](%s)", title, url)

	if notifType == "" {
		return fmt.Sprintf("* %s %s\n", repoPart, titlePart)
	}

	return fmt.Sprintf("* %s %s %s\n", repoPart, notifType, titlePart)
}

func fullNameFromOwnerAndRepo(owner, repo string) string {
	return fmt.Sprintf("%s/%s", owner, repo)
}
