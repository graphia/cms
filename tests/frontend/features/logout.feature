Feature: Logging in
	So I can securely stop using the application
	As an editor
	I want to log out of the CMS

	Background:
		Given the CMS is running with the "default" config
		And a user account has been created

	Scenario: Logging out
		Given I am logged in
		When I select 'Logout' from the settings menu
		Then I should be logged out
		And I should be on the login screen
		And there should be no entries in the navigation bar