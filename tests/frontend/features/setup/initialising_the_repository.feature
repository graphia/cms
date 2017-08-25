Feature: Setting up an initial user

	So I have a workspace for my documents
	As an customer
	I want to set up a workspace

	Background:
		Given a user account has been created
		And I have logged in
		And there is an empty directory in place of a repository

	Scenario: Redirecting to the initialise repo screen
		Given I try to navigate to the home page
		Then I should be redirected to the initialize repository page

	Scenario: Initialize repository page contents
		Given I am on the initialize repository page
		Then I should see a button with text 'Initialise Repository'