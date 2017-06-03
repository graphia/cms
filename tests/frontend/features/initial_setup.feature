Feature: First run

	So I can configure the system for use
	As an customer
	I want to create the first user account

	Scenario: Creating an initial user
		Given there are no users
		When I navigate to the login page
		Then I should be redirected to the initial setup page
