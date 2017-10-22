Feature: Deleting documents
	So I can remove unwanted content to the CMS
	As an author
	I want to be able to delete documents

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: Deleting a file
		Given I am on the document's show page
		When I click the "Delete" button
		Then I should be redirected to the parent directory's index
		And the file should have been deleted

	Scenario: Commit message
		Given I have deleted a single file
		Then the last commit message should contain the file's name

	Scenario: Should show a appropriate error when repo out of sync
		Given I am on the document's show page
		And a repository update has taken place in the background
		When I click the "Delete" button
		Then there should be an alert with the message "The repository is out of sync"

	Scenario: Deleting a file after reloading data
		Given I have tried to delete a file after a repo update
		When I click the "Delete" button again
		Then I should be redirected to the parent directory's index
		And the file should have been deleted