Feature: Git via SSH
	So I can work on content offline
	As a power user
	I want to be able to clone the repository via SSH

	Background:
		Given a repository has been initialised
		And the CMS is running with the "ssh_enabled" config
		And my user account with public key exists

	Scenario: Connecting to the server with a valid key
		Given my private key is valid
		When I initiate a SSH connection to the server
		Then I should see the response "Graphia: Connection successful"

	Scenario: Connecting to the server with an invalid key
		Given my private key is invalid
		When I initiate a SSH connection to the server
		Then I should receive an AuthenticationFailed error

	Scenario: Triyng to connect with a user other than git
		Given my private key is valid
		When I try to establish a connection with user "krusty"
		Then I should receive the error message "Access denied"

	Scenario: Trying to run an illegal command
		Given my private key is valid
		When I try to run one of the following commands:
			| ls -la |
			| cd /   |
			| bash   |
		Then I should receive an error

	Scenario: Cloning the content repository
		Given I have an SSH key
		When I try to clone the "content" repository
		Then the directory should be present in my working directory
		And I should see output detailing my clone operation

	Scenario: Attempting to clone a non-existant repository
		Given my private key is valid
		When I try to clone the "non_existant" repository
		Then I should see output with an error

	Scenario: Cloning the theme repository
		Given I have an SSH key
		And a theme is present
		When I try to clone the "theme" repository
		Then the directory should be present in my working directory
		And I should see output detailing my clone operation