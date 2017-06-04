Feature: First run

	So I can configure the system for use
	As an customer
	I want to create the first user account

	Scenario: Creating an initial user
		Given there are no users
		When I navigate to the login page
		Then I should be redirected to the initial setup page

	Scenario: The initial setup form contents
		Given I am on the initial setup page
		Then I should see a form with the following fields:
			| Name             | Type     | Required  |
			| Name             | Text     | yes       |
			| Username         | Text     | yes       |
			| Password         | Password | yes       |
			| Confirm Password | Password | yes       |

		And the submit button should be labelled 'Create initial user'

	Scenario: Entering a username that is too short
		Given I am on the initial setup page
		When I enter a '3' letter word into 'Username'
		Then the 'Username' field should be invalid

	Scenario: Entering an valid length username
		Given I am on the initial setup page
		When I enter a '6' letter word into 'Username'
		Then the 'Username' field should be valid

	Scenario: Non-matching password confirmation
		Given I am on the initial setup page
		When I enter 'p4s5w0rD' in the 'Password' field
		And I enter 's3cRetz' in the 'Confirm Password' field
		Then the 'Confirm Password' field should be marked invalid
		And there should be a warning containing 'Password and confirmation must match'

	Scenario: When the passwords do match
		Given I am on the initial setup page
		When I enter matching passwords in the 'Password' and 'Confirm Password' fields
		Then no password-related warnings should be visible