Feature: Logging in

	So I can manage content
	As an editor
	I want to log into the CMS

	Scenario: Login screen contents
		Given I am on the login screen
		Then I should see a 'Username' field with type 'text'
		Then I should see a 'Password' field with type 'password'
		And the submit button should be labelled 'Log in'

	Scenario: Logging in with invalid credentials
		Given I am on the login screen
		When I enter invalid credentials
		And I submit the form
		Then I should still be on the login screen
		And there should be a 'red' alert with the message 'Invalid'

	Scenario: Creating an initial user
		Given there are no users
		When I navigate to the CMS
		Then I should see the initial user setup form

