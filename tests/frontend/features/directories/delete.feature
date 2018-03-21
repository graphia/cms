Feature: Creating directories
	So I can organise documents
	As an author
	I want to be able to remove unwanted directories

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Delete button presence
		Given I am on the "documents" index page
		Then I should see a "Delete directory" button

	Scenario: The deletion modal
		Given I am on the "documents" index page
		When I click the "Delete directory" button
		Then I should see the deletion modal box

	Scenario: Actually deleting a directory
		Given I am on the "documents" index page
		When I click 'Delete directory' then 'Confirm deletion'
		Then I should see a message containing 'Documents and its contents have been deleted'
		And I should be on the homepage
		And the "documents" directory should have been deleted