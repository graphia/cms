Feature: Deleting documents
	So I can remove unwanted content to the CMS
	As an author
	I want to be able to delete documents

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: The deletion modal
		Given I am on the document's show page
		When I click the "Delete" button
		Then I should see the deletion modal box

	Scenario: Deleting a file
		Given I can see the document's deletion modal
		When I click the "Confirm deletion" button
		Then I should be redirected to the parent directory's index
		And the file should have been deleted

	Scenario: Commit message
		Given I have deleted a single file
		Then the last commit message should contain the file's name

	Scenario: Deleting a file when the repository has been updated in the background
		Given I am on the document's show page
		And a repository update has taken place in the background
		When I try to delete the file
		Then I should be redirected to the parent directory's index
		And the file should have been deleted

	Scenario: Deleting a file plus its attachments
		Given I can see the document's deletion modal
		When I check the "Delete attachments" checkbox
		And I click the "Confirm deletion" button
		Then I should be redirected to the parent directory's index
		And the file and attachments directory should have been deleted

	Scenario: Deleting a file without its attachments
		Given I can see the document's deletion modal
		When I click the "Confirm deletion" button
		Then I should be redirected to the parent directory's index
		And the file should have been deleted but not the attachments directory