Feature: The commit message
	So I can't submit the form without sufficient information
	As an author
	I want the submit button to enable itself when the data is valid

	Background:
		Given my user account exists
		And I have logged in
		And I am on the new document page

	Scenario: Submit disabled when form is initialised
		Given I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: Submit button is disabled when form is invalid
		Given I enter invalid information in the form
		Then the submit button should be disabled

	Scenario: Submit button is enabled when form is valid
		Given I enter valid information in the form
		Then the submit button should be enabled
