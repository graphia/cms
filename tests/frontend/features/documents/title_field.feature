Feature: The title field
	So I can organise and format my documents properly
	As an author
	I want to be able to give them custom titles

	Background:
		Given my user account exists
		And I have logged in
		And I am on the new document page

	Scenario: Validation message when field is untouched
		Given I haven't touched the 'Title' field
		Then the title validation feedback should not be visible

	Scenario: Validation message appears when the field is edited
		Given I enter 'a' in the 'Title' field
		Then the title validation feedback should be visible

	Scenario: Validation message is not displayed when the value is valid
		Given I enter valid information in the form
		Then the title validation feedback should not be visible