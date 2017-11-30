Feature: Setting up my SSH Keys
	So I can work on content offline
	As a power user
	I want to be able to configure my SSH Public Keys

	Background:
		Given a repository has been initialised
		And the CMS is running with the "ssh_enabled" config
		And my user account exists
		And I have logged in

	Scenario: Settings page contents
		Given I am on my ssh key settings page
		Then I should see the heading "User Settings"
		And I should see subheadings:
			 | Existing keys        |
			 | Upload a new SSH key |

	Scenario: When I have no SSH keys I should see an informative message
		Given I have no SSH keys
		When I visit my ssh key settings page
		Then I should see an alert informing me why I might want to add one

	Scenario: When I have SSH keys I should see them in a list
		Given I have an SSH key
		When I visit my ssh key settings page
		Then I should see my keys listed by name and fingerprint
		And each entry should have a 'Delete' button

	Scenario: New SSH key form
		Given I am on my ssh key settings page
		Then I should see a 'name' field with type 'text'
		And I should see a text area called 'ssh-key'
		And the submit button should be labelled 'Create SSH Key'

	Scenario: Adding an SSH key
		Given I am on my ssh key settings page
		And I enter 'home-desktop' in the 'Name' field
		And I paste my public SSH key into the 'ssh-key' text area
		When I submit the form
		Then there should be an alert with the message "SSH Key Created"
		And my SSH Key should be added to the Existing Keys list

	Scenario: Deleting an SSH key
		Given I have an SSH key
		And I visit my ssh key settings page
		When I click the 'Delete' button for my SSH key
		Then there should be an alert with the message "SSH Key Deleted"
		And my key should have been deleted