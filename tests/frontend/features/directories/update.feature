Feature: Creating directories
	So I can keep information up to date
	As an author
	I want to be able to modify directory metadata

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in


	Scenario: The edit directory page
		Given I am on the update directory page
		Then I should see a form with the following fields:
			| Name  | Type | Required |
			| Title | Text | yes      |
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
		And the submit button should be labelled 'Update directory'

	Scenario: Page title
		Given I am on the update directory page
		And the page's title should be "Edit directory"

	Scenario: File updates
		Given I am on the update directory page
		And I fill in the form with the following data:
			| Title       | Ice Creams                                  |
			| Description | A description of ice cream related products |
		And I set the editor text to "# Fabulous ices of all colours"
		When I submit the form
		Then I am on the "documents" index page
		And the directory index page should contain the newly-updated information