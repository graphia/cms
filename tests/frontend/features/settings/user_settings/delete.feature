Feature: User creation

	So I can add colleagues to the system
	As a administrator
	I want to be able to create new user accounts

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Deleting a user
		Given there is a regular user and an administrator
		And I am on the edit user page for the regular user
		When I click the "Delete" button
		Then I should be on the users list page
		And the regular user should not be present
		And I should see a message containing 'User Herman Hermann deleted'

	Scenario: Not being able to delete youself
		Given I am on my own edit user page
		When I click the "Delete" button
		Then I should see an error containing "not deleted"