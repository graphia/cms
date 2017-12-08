Feature: Settings page

	So I can keep my account secure
	As a user
	I want to be able to modify my password

	Background:
		Given the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page contents
		Given I am on my settings page
		Then I should see a section called "Change my password"
		And I should see a "Passwords" form with the following fields:
			| Name             | Type     | Required | Disabled |
			| Current Password | Password | Yes      | No       |
			| Confirm Password | Password | Yes      | No       |
			| New Password     | Password | Yes      | No       |

	Scenario: Non-matching confirmation
		Given I am on my settings page
		When I enter 'abc123' in the 'New Password' field
		And I enter 'def456' in the 'Confirm Password' field
		Then the "Confirm Password" field should be marked invalid

	Scenario: Non-matching confirmation
		Given I am on my settings page
		When I enter 'abc123' in the 'New Password' field
		And I enter 'abc123' in the 'Confirm Password' field
		Then the "New Password" and "Confirm Password" fields should be marked valid

	Scenario: Attempting update with correct current password
		Given I am on my settings page
		And I enter 'okily-dokily!' in the 'Current Password' field
		And I enter 'neighbourino' in the 'New Password' field
		And I enter 'neighbourino' in the 'Confirm Password' field
		When I submit the "Passwords" form
		Then I should see a message containing 'Password updated'

	Scenario: Attempting update with incorrect current password
		Given I am on my settings page
		And I enter 'toodily-doo' in the 'Current Password' field
		And I enter 'neighbourino' in the 'New Password' field
		And I enter 'neighbourino' in the 'Confirm Password' field
		When I submit the "Passwords" form
		Then I should see an error containing "Current password is not correct"