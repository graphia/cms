Feature: Setting up an initial user

	So I can configure the system for use
	As an customer
	I want to create the first user account

	Background:
		Given the CMS is running with the "default" config
		And there are no users

	Scenario: Creating an initial user
		Given I navigate to the login page
		Then I should be redirected to the initial setup page

	Scenario: The initial setup form contents
		Given I am on the initial setup page
		Then I should see a form with the following fields:
			| Name             | Type     | Required  |
			| Full Name        | Text     | yes       |
			| Username         | Text     | yes       |
			| Email Address    | Email    | yes       |
			| Password         | Password | yes       |
			| Confirm Password | Password | yes       |
		And the submit button should be labelled 'Create'

	Scenario: HTML5 input length validation attributes
		Given I am on the initial setup page
		Then the 'Username' field should allow values from '3' to '32' characters
		And the 'Full Name' field should allow values from '3' to '64' characters
		And the 'Password' field should be at least '6' characters long

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

	Scenario: Creating an administrator
		Given I am on the initial setup page
		And I fill in the form with the following data:
			| Full Name        | Patty Bouvier     |
			| Username         | p.bouvier         |
			| Email Address    | patty.b@yahoo.com |
			| Password         | macguyver101      |
			| Confirm Password | macguyver101      |
		When I submit the form
		Then I should see a message containing 'Administrator created'
		And the new user should have been saved to the database