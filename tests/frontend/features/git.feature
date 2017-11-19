Feature: Commits
	So I can work on content offline
	As a power user
	I want to be able to clone the repository via SSH

	Background:
		Given a repository has been initialised
		And the CMS is running with the "ssh_enabled" config
		And my user account with public key exists

	@wip
	Scenario: Connecting to the server with a valid key
		Given my private key is valid
		When I initiate a SSH connection to the server
		Then I should see the response "Graphia: Connection successful"

	Scenario: Connecting to the server with an invalid key
		Given my private key is invalid
		When I initiate a SSH connection to the server
		Then I should receive an AuthenticationFailed error