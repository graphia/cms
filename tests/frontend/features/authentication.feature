Feature: Authentication
	So the content of the CMS is protected
	As a user
	I want be redirected to the login screen if I'm not authenticated

	Background:
		Given my user account exists

	Scenario: Navigating to the homepage with no session
		Given I am not logged in
		When I navigate directly to the homepage
		Then I should be redirected to the login page

	Scenario: Directly accessing the API when not logged in
		Given I am not logged in
		When I try to manually make an unauthenticated HTTP request to the API
		Then I should receive a 'HTTP Not Authorized' status

	Scenario: Directly accessing the API when logged in
		Given I have logged in
		When I try to manually make an authenticated HTTP request to the API
		Then I should receive a 'HTTP OK' status