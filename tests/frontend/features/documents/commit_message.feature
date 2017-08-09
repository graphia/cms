Feature: The commit message
	So I can add information about changes to a commit
	As an author
	I want to be able to create a valid commit message

	Background:
		Given my user account exists
		And I have logged in
		And I am on the new document page

	Scenario: Validation message when field is untouched
		Given I haven't touched the 'Commit Message' field
		Then the commit message validation feedback should not be visible

	Scenario: Validation message appears when the field is edited
		Given I enter 'abc' in the 'Commit Message' field
		Then the commit message validation feedback should be visible

	Scenario: Validation message is not displayed when the value is valid
		Given I enter valid information in the form
		Then the commit message validation feedback should not be visible