Feature: Logging in

	So I can manage content
	As an editor
	I want to log into the CMS

	Scenario: Login screen contents
		Given I am on the login screen
		Then I should see a 'Username' field with type 'text'
		Then I should see a 'Password' field with type 'password'
		And the should be a submit button should be labelled 'Login'
