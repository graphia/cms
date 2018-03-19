Feature: Setting up an initial user

	So I can begin managing content
	As a newly-created user
	I want to set my password and login

	Background:
		Given the CMS is running with the "default" config
		And there is a user requiring activation

	Scenario: The user activation screen
		Given I am on the activate user screen
		Then I should see a personal welcome
		And there should be a password reset form
		And the password fields should be required and at least 6 characters long

	Scenario: Activating a user
		Given I am on the activate user screen
		And I fill in and confirm the password
		When I submit the form
		Then I should see a message containing 'Your account has been activated'
		And I should be on the login screen

	Scenario: Logging in with an activated user
		Given I have activated my account
		And I am on the login screen
		When I enter my username and password
		And I submit the form
		Then I should be logged in