Feature: Settings page

	So I can tailor the system to my needs
	As a user
	I want to be able to view and edit my personal details

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in


	Scenario: Page contents
		Given I am on my settings page
		Then the page heading should be "Settings"
		And I should see subheadings:
			| Personal Details |