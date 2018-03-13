Feature: User creation

	So I can add colleagues to the system
	As a administrator
	I want to be able to create new user accounts

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page contents
		Given I am on the edit user page
		Then the page's title should be "Edit user"
		And I should see a form with the following fields:
			| Name          | Type        | Required |
			| Name          | Text        | yes      |
			| Email address | Email       | yes      |
			| Administrator | Checkbox    | no       |

	Scenario: Validation
		Given I am on the edit user page
		And the 'Name' field should allow values from '3' to '64' characters

	Scenario: Editing a user
		Given I am on the edit user page
		When I fill in the form with the following data:
			| Name             | Todd Flanders |
			| Email address    | tf@floody.com |
		And I submit the form
		Then I should see a message containing 'User Todd Flanders updated'
		And I see my newly-updated user when redirected to the user list

	Scenario: Cancelling a user edit
		Given I am on the edit user page
		When I click the "Cancel" button
		Then I should be on the users list page

	Scenario: Server-side validation errors
		Given there is a regular user and an administrator
		And I am on the edit user page
		When I change the name and email address to match the regular user
		And I submit the form
		Then I should see an error message stating that the record is not unique

	Scenario: Breadcrumbs
		Given I am on the edit user page
		Then I should see the following breadcrumbs:
			| Text      | Reference           |
			| Dashboard | /cms                |
			| Users     | /cms/settings/users |
			| Edit      | None                |

	Scenario: Non-existing user
		Given I am on the edit page for non-existing user 'julius.hibbert'
		Then I should see "404"
		And I should see "User julius.hibbert cannot be found"
