Feature: Logging in

	So I can manage content
	As an editor
	I want to log into the CMS

	Background:
		Given a user account has been created

	Scenario: Login screen contents
		Given I am on the login screen
		Then I should see a form with the following fields:
			| Name             | Type     | Required  |
			| Username         | Text     | yes       |
			| Password         | Password | yes       |
		And the submit button should be labelled 'Log in'

	Scenario: Logging in with invalid credentials
		Given I am on the login screen
		When I enter invalid credentials
		And I submit the form
		Then I should still be on the login screen
		And there should be an alert with the message "Invalid credentials"

	Scenario: Logging in with valid credentials
		Given I am on the login screen
		When I enter valid credentials
		And I submit the form
		Then I should see a message containing 'You have logged in successfully'
		And I should be redirected to the CMS's landing page

	Scenario: Logging in with valid credentials
		Given I have logged in
		Then I should have a JWT saved in localstorage