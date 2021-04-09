package xErr

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHandleErrorTypes(t *testing.T) {
	type unit struct {
		givenError       error
		when             string
		then             string
		expectedExitCode int
		expectedMessage  string
	}

	test := []unit{
		{
			givenError:       ErrorAlreadyInLastBranch,
			when:             "When ErrorAlreadyInLastBranch error",
			then:             "Then break execution with exit code 0 and show proper message",
			expectedExitCode: 0,
			expectedMessage:  "Error:\n  already on the last branch of the deployment workflow, no need to send further PRs\n",
		},
		{
			givenError:       ErrorOnlyOneBranchInWorkflow,
			when:             "When ErrorOnlyOneBranchInWorkflow error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  only one branch in workflow arguments, minimum two\n",
		},
		{
			givenError:       ErrorTravisRepoSlug,
			when:             "When ErrorTravisRepoSlug error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  TRAVIS_REPO_SLUG is not defined\n",
		},
		{
			givenError:       ErrorTravisBuildNumber,
			when:             "When ErrorTravisBuildNumber error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  TRAVIS_BUILD_NUMBER is not defined\n",
		},
		{
			givenError:       ErrorTravisBuildID,
			when:             "When ErrorTravisBuildID error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  TRAVIS_BUILD_ID is not defined\n",
		},
		{
			givenError:       ErrorTravisBranch,
			when:             "When ErrorTravisBranch error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  TRAVIS_BRANCH is not defined\n",
		},
		{
			givenError:       ErrorTravisIsAPR,
			when:             "When ErrorTravisIsAPR error",
			then:             "Then break execution with exit code 0 and show proper message",
			expectedExitCode: 0,
			expectedMessage:  "Error:\n  it's a PR, won't launch go-pr-creator\n",
		},
		{
			givenError:       ErrorGitHubToken,
			when:             "When ErrorGitHubToken error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  GITHUB_TOKEN is not defined\n",
		},
		{
			givenError:       ErrorNoArgs,
			when:             "When ErrorNoArgs error",
			then:             "Then break execution with exit code 0 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  no arguments given\n",
		},
		{
			givenError:       ErrorWorkflowBranch,
			when:             "When ErrorWorkflowBranch error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  workflow or branch params must be given\n",
		},
		{
			givenError:       ErrorConnectionToGithub,
			when:             "When ErrorConnectionToGithub error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  error connecting to github\n",
		},
		{
			givenError:       ErrorPullRequest,
			when:             "When ErrorPullRequest error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  error creating Pull request\n",
		},
		{
			givenError:       ErrorLabels,
			when:             "When ErrorLabels error",
			then:             "Then DON'T break execution and show proper message",
			expectedExitCode: -1,
			expectedMessage:  "Error:\n  error adding labels to the PR\n",
		},
		{
			givenError:       ErrorGitHubClient,
			when:             "When errorGitHubClient error",
			then:             "Then break execution with exit code 1 and show proper message",
			expectedExitCode: 1,
			expectedMessage:  "Error:\n  error creating github client\n",
		},
		{
			givenError:       ErrorNot200,
			when:             "When errorNot200 error",
			then:             "Then DON'T break execution and show proper message",
			expectedExitCode: -1,
			expectedMessage:  "Error:\n  got an http code different to 200\n",
		},
		{
			givenError:       ErrorListingCommits,
			when:             "When errorListingCommits error",
			then:             "Then DON'T break execution and show proper message",
			expectedExitCode: -1,
			expectedMessage:  "Error:\n  error obtaining the list of commits in this pull request. Please, review them manually\n",
		},
		{
			givenError:       ErrorUpdatingBody,
			when:             "When errorUpdatingBody error",
			then:             "Then DON'T break execution and show proper message",
			expectedExitCode: -1,
			expectedMessage:  "Error:\n  error updating body message\n",
		},
	}

	// Mock os.Exit
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()
	var exitCode int
	myExit := func(code int) {
		exitCode = code
	}
	osExit = myExit

	oldOutput := output
	defer func() { output = oldOutput }()

	Convey("Given the parameters sent to the go-pr-creator", t, func() {
		for _, unit := range test {
			Convey(unit.when, func() {
				exitCode = -1
				output = new(bytes.Buffer)
				HandleError(unit.givenError)
				Convey(unit.then, func() {
					So(exitCode, ShouldEqual, unit.expectedExitCode)
					So(output.(*bytes.Buffer).String(), ShouldEqual, unit.expectedMessage)
				})
			})
		}
	})
}

func TestHandleErrorWrapped(t *testing.T) {
	oldOutput := output
	defer func() { output = oldOutput }()

	Convey("Given an error", t, func() {
		Convey("When the error is wrapped", func() {
			exitCode := -1
			output = new(bytes.Buffer)
			err := errors.New("whatever error")
			wrap := fmt.Errorf("%w:\n    %v", err, "Test message")
			HandleError(wrap)
			fmt.Println(wrap)

			Convey("We get true and the expected index in the array", func() {
				So(exitCode, ShouldEqual, -1)
				So(output.(*bytes.Buffer).String(), ShouldEqual, "Error:\n  whatever error:\n    Test message\n")
			})
		})
	})
}

func TestWrapError(t *testing.T) {
	Convey("Given two errors", t, func() {
		Convey("When we wrap them", func() {
			err := errors.New("whatever error")
			wrap := errors.New("test message")
			obtained := WrapError(err, wrap)
			expected := errors.New("whatever error:\n    test message")

			Convey("We wrap them with the correct formatting", func() {
				So(obtained.Error(), ShouldResemble, expected.Error())
			})
		})
	})
}

func TestUnwrapErrorNoNested(t *testing.T) {
	Convey("Given an error that hasn't been wrapped", t, func() {
		Convey("When we unwrap it", func() {
			err := errors.New("error id")
			obtained := UnwrapError(err)
			expected := err

			Convey("We obtain the same error", func() {
				So(obtained.Error(), ShouldResemble, expected.Error())
			})
		})
	})
}

func TestUnwrapErrorNested(t *testing.T) {
	Convey("Given an error that has been wrapped", t, func() {
		Convey("When we unwrap it", func() {
			err := errors.New("error id")
			wrap := fmt.Errorf("wrapping message: %w", err)
			obtained := UnwrapError(wrap)
			expected := err

			Convey("We obtain the first error", func() {
				So(obtained.Error(), ShouldResemble, expected.Error())
			})
		})
	})
}
