Feature: Setting up an initial user

	So I have a workspace for my documents
	As an customer
	I want to set up a workspace

	Background:
		Given there are no users
		When I navigate to the login page
		Then I should be redirected to the initial setup page