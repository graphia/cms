Feature: Home page contents

	So I can see an overview of the CMS
	As a user
	I want the homepage to provide me with useful information

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: Home page sections
		Given I am on the homepage
		Then I should see a summary of recent changes
		And I should see a statistics section

	Scenario: Recent commits
		Given there have been some recent commits
		And I am on the homepage
		Then the recent changes summary should contain a list of commits