Feature: Commits
	So I can review historical changes
	As a user
	I want to be able to view individual commits and their details

	Background:
		Given my user account exists
		And I have logged in

	Scenario: Commit metadata
		Given I have added a new file
		When I navigate to the commit's details page
		Then I should see some information specific to the commit

	Scenario: Added files
		Given I have added a new file
		When I navigate to the commit's details page
		Then I should see a 'green' section with the file's path for a title
		And it should contain the entire file's contents

	Scenario: Modified files
		Given I have made changes to an existing file
		When I navigate to the commit's details page
		Then I should see a 'blue' section with the file's path for a title
		And it should contain a colourised diff showing changes made

	Scenario: Deleted files
		Given I have deleted a file
		When I navigate to the commit's details page
		Then I should see a 'red' section with the file's path for a title
		And it should contain the entire file's contents