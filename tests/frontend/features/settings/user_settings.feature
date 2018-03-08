Feature: User settings

	So I can view a list of colleagues
	As a user
	I want to be able to see a list of registered users

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page contents
		Given I am on the users list page
		Then there should be a user list
		And there should be a "Create New User" button

	Scenario: User list
		Given that '2' extra users have been created
		When I am on the users list page
		Then I should see a user list with '3' entries

	Scenario: List item contents
		Given I am on the users list page
		Then I should see a section with my user's name as the title
		And the details listed should be:
			| Username | rod.flanders                               |
			| Email    | rod.flanders@springfield.elementary.k12.us |

	Scenario: Navigating to the edit page
		Given I am on the users list page
		When I click the "Edit" button for my user
		Then I should be on the edit page for my user

	Scenario: Navigating to the new user page
		Given I am on the users list page
		When I click the "Create New User" link
		Then I should be on the new user page