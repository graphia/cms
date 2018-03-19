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
		And I should see a "Personal Details" form with the following fields:
			| Name     | Type  | Required | Disabled |
			| Username | Text  | Yes      | Yes      |
			| Email    | Email | Yes      | Yes      |
			| Name     | Text  | Yes      | No       |
		And the "Personal Details" submit button should be disabled

	Scenario: Page title
		Given I am on my settings page
		Then the page's title should be "Settings: Personal Details"

	Scenario: Enabling of the submit button
		Given I am on my settings page
		When I change my name to "Ranier Wolfcastle"
		Then the "Personal Details" submit button should be enabled

	Scenario: Updating my name
		Given I am on my settings page
		When I change my name to "Ranier Wolfcastle"
		And I submit the "Personal Details" form
		Then my name should have changed to "Ranier Wolfcastle"

	Scenario: Breadcrumbs
		Given I am on my settings page
		Then I should see the following breadcrumbs:
			| Text                | Reference   |
			| Dashboard           | /cms        |
			| My Profile          | None        |
