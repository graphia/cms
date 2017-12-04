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
		Then the page heading should be "Personal Details"
		And I should see a form with the following fields:
			| Name     | Type  | Required | Disabled |
			| Username | Text  | Yes      | Yes      |
			| Email    | Email | Yes      | Yes      |
			| Name     | Text  | Yes      | No       |
		And the submit button should be disabled

	Scenario: Enabling of the submit button
		Given I am on my settings page
		When I change my name to "Ranier Wolfcastle"
		Then the submit button should be enabled

	Scenario: Updating my name
		Given I am on my settings page
		When I change my name to "Ranier Wolfcastle"
		And I submit the form
		Then my name should have changed to "Ranier Wolfcastle"