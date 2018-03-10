Feature: User creation

	So I can add colleagues to the system
	As a administrator
	I want to be able to create new user accounts

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page contents
		Given I am on the new user page
		Then the page's title should be "New user"
		And I should see a form with the following fields:
			| Name          | Type        | Required |
			| Username      | Text        | yes      |
			| Name          | Text        | yes      |
			| Email address | Email       | yes      |
			| Administrator | Checkbox    | no       |

	Scenario: Validation
		Given I am on the new user page
		Then the 'Username' field should allow values from '3' to '32' characters
		And the 'Name' field should allow values from '3' to '64' characters

	Scenario: Creating a new user
		Given I am on the new user page
		When I fill in the form with the following data:
			| Username         | hhermann       |
			| Name             | Herman Hermann |
			| Email address    | hello@hma.com  |
		And I submit the form
		Then I should see a message containing 'Herman Hermann will receive an email with instructions on how to log in'
		And I see my newly-created user when redirected to the user list

	Scenario: Cancelling a new user creation
		Given I am on the new user page
		When I click the "Cancel" button
		Then I should be on the users list page

	Scenario: Server-side validation errors
		Given I am on the new user page
		And I re-enter the details of an existing user
		When I submit the form
		Then I should see an error message