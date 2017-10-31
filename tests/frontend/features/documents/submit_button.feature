Feature: The submit button
	So I can't submit the form without sufficient information
	As an author
	I want the submit button to enable itself when the data is valid

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Submit disabled when the create document form is initialised
		Given I am on the new document page
		When I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: Submit button is disabled when the create document form is invalid
		Given I am on the new document page
		When I enter invalid information in the form
		Then the submit button should be disabled

	Scenario: Submit button is enabled when the create document form is valid
		Given I am on the new document page
		When I enter valid information in the form
		Then the submit button should be enabled

	Scenario: Submit disabled when the update document form is initialised
		Given I am on the edit document page for a document
		When I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: Submit button is disabled when the update document form is invalid
		Given I am on the edit document page for a document
		When I enter invalid information in the form
		Then the submit button should be disabled

	Scenario: Submit button is enabled when the update document form is valid
		Given I am on the edit document page for a document
		When I enter valid information in the form
		Then the submit button should be enabled
