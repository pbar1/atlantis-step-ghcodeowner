// Copyright (c) 2021 Pierce Bartine. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package atlantis

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (

	// RunStep is a struct populated with the current Atlantis custom run command runtime variables.
	RunStep struct {

		// The Terraform workspace used for this project, ex. `default`.
		// NOTE: if the step is executed before `init` then Atlantis won't have switched to this workspace yet.
		TerraformWorkspace string `envconfig:"WORKSPACE" required:"true"`

		// The version of Terraform used for this project, ex. `0.11.0`.
		TerraformVersion string `envconfig:"ATLANTIS_TERRAFORM_VERSION" required:"true"`

		// Absolute path to the current directory.
		AbsoluteProjectDir string `envconfig:"DIR" required:"true"`

		// Absolute path to the location where Atlantis expects the plan to either be generated (by plan) or already exist (if running apply). Can be used to override the built-in `plan`/`apply` commands, ex. run: `terraform plan -out $PLANFILE`.
		Planfile string `envconfig:"PLANFILE" required:"true"`

		// Name of the repository that the pull request will be merged into, ex. `atlantis`.
		BaseRepoName string `envconfig:"BASE_REPO_NAME" required:"true"`

		// Owner of the repository that the pull request will be merged into, ex. `runatlantis`.
		BaseRepoOwner string `envconfig:"BASE_REPO_OWNER" required:"true"`

		// Name of the repository that is getting merged into the base repository, ex. `atlantis`.
		HeadRepoName string `envconfig:"HEAD_REPO_NAME" required:"true"`

		// Owner of the repository that is getting merged into the base repository, ex. `acme-corp`.
		HeadRepoOwner string `envconfig:"HEAD_REPO_OWNER" required:"true"`

		// Name of the head branch of the pull request (the branch that is getting merged into the base)
		HeadBranchName string `envconfig:"HEAD_BRANCH_NAME" required:"true"`

		//  Name of the base branch of the pull request (the branch that the pull request is getting merged into)
		BaseBranchName string `envconfig:"BASE_BRANCH_NAME" required:"true"`

		// Name of the project configured in `atlantis.yaml`. If no project name is configured this will be an empty string.
		ProjectName string `envconfig:"PROJECT_NAME"`

		// Pull request number or ID, ex. `2`.
		PullNum int `envconfig:"PULL_NUM" required:"true"`

		// Username of the pull request author, ex. `acme-user`.
		PullAuthor string `envconfig:"PULL_AUTHOR" required:"true"`

		// The relative path of the project in the repository. For example if your project is in `dir1/dir2/` then this will be set to `dir1/dir2`. If your project is at the root this will be `"."`.
		RelativeProjectDir string `envconfig:"REPO_REL_DIR" required:"true"`

		// Username of the VCS user running command, ex. acme-user. During an `autoplan`, the user will be the Atlantis API user, ex. `atlantis`.
		Username string `envconfig:"USER_NAME" required:"true"`

		// Any additional flags passed in the comment on the pull request.
		CommentArgs CommentArgs `envconfig:"COMMENT_ARGS"`
	}

	// CommentArgs is a string slice.
	CommentArgs []string
)

// Decode parses the Atlantis `COMMENT_ARGS` environment variable into a string slice.
// Flags are separated by commas and every character is escaped, ex. `atlantis plan -- arg1 arg2` will result in `COMMENT_ARGS=\a\r\g\1,\a\r\g\2`.
func (a *CommentArgs) Decode(value string) error {
	args := make([]string, 0)

	if value == "" {
		*a = args
		return nil
	}

	if strings.Count(value, `\,`)*2 == len(value) {
		return fmt.Errorf("ambiguous string: %s", value)
	}

	// append trailing comma to make logic simpler
	escaped := []rune(value + ",")

	var arg []rune
	lastCommaPos := -1
	for i, c := range escaped {
		// if we reach a comma that's not escaped (ie, an actual arg separator)
		if c == ',' && (escaped[i-1] != '\\' || escaped[i-1] == '\\' && escaped[i-2] == '\\') {
			// take the substring between the last separator and this separator (ie, the escaped arg, minus the comma)
			arg = escaped[lastCommaPos+1 : i]
			// add the unescaped arg to the return list
			uarg, err := unescape(arg)
			if err != nil {
				return fmt.Errorf(`unable to unescape "%s": %v`, value, err)
			}
			args = append(args, uarg)
			// update last separator index to current separator index (ie, anchor the beginning of the next arg)
			lastCommaPos = i
		}
	}

	*a = args
	return nil
}

// unescape takes every other character of a rune slice and returns the resulting string. In Atlantis's case, every
// character is escaped by a backslash, but it could as well be any escape character.
func unescape(runes []rune) (string, error) {
	unescaped := []rune("")
	for i, c := range runes {
		if i%2 == 0 && c != '\\' {
			return "", fmt.Errorf("improperly escaped arg: %s", string(runes))
		}
		if i%2 == 1 {
			unescaped = append(unescaped, c)
		}
	}
	return string(unescaped), nil
}

// NewRunStep constructs a RunStep populated with the current Atlantis custom run command runtime variables.
func NewRunStep() (*RunStep, error) {
	var runStep RunStep
	if err := envconfig.Process("", &runStep); err != nil {
		return nil, fmt.Errorf("unable to populate RunStep: %v", err)
	}
	return &runStep, nil
}
