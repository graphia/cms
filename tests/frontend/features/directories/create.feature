Feature: Creating directories
	So I can organise documents
	As an author
	I want to be able to create new directories

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: The new directory page
		Given I am on the new directory page
		Then I should see a form with the following fields:
			| Name  | Type | Required |
			| Title | Text | yes      |
			| Path  | Text | yes      |
		And I should see a text area called 'Description'
		And I should see an editor with the following buttons:
			| Bold           |
			| Italic         |
			| Heading        |
			| Quote          |
			| Generic List   |
			| Numbered List  |
			| Create Link    |
			| Insert Image   |
			| Toggle Preview |
			| Markdown Guide |
		And the submit button should be labelled 'Create directory'

	Scenario: Creating a directory
		Given I am on the new directory page
		And I fill in the form with the following data:
			| Title       | Ice Creams                                  |
			| Path        | ice-creams                                  |
			| Description | A description of ice cream related products |
		And I set the editor text to "# Fabulous ices of all colours"
		When I submit the form
		Then I should be redirected to the new directroy's index
		And the directory should have been created with the correct information