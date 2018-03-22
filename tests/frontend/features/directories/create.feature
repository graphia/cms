Feature: Creating directories
	So I can organise documents
	As an author
	I want to be able to create new directories

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
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

	Scenario: Page title
		Given I am on the new directory page
		And the page's title should be "Create directory"

	Scenario: Creating a directory
		Given I am on the new directory page
		And I fill in the form with the following data:
			| Title       | Ice Creams                                  |
			| Description | A description of ice cream related products |
		And I set the editor text to "# Fabulous ices of all colours"
		When I submit the form
		Then I should be redirected to the new directroy's index
		And the directory should have been created with the correct information

	Scenario: The path field should be the slugged version of the title
		Given I am on the new directory page
		When I enter "Kwik E Mart" in the "Title" field
		Then the "Path" field should be "kwik-e-mart"

	Scenario: The submit button shouldn't be active by default
		Given I am on the new directory page
		When I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: The submit button should become active when the form is valid
		Given I am on the new directory page
		When I have filled in the required fields
		Then the submit button should be enabled

	Scenario: Cancelling a directory creation
		Given I am on the new directory page
		When I click the "Cancel" button
		Then I should be on the homepage

	Scenario: Displaying error messages
		Given I am on the new directory page
		When I enter "f" in the "Title" field
		Then the submit button should be disabled
		And the "Title" field's error text should contain "Please lengthen this text to 2 characters or more"

	Scenario: Breadcrumbs
		Given I am on the new directory page
		Then I should see the following breadcrumbs:
			| Text                | Reference       |
			| Dashboard           | /cms            |
			| New Directory       | None            |