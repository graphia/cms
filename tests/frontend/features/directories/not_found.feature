Feature: Displaying the correct message when a directory is not found
	So I'm informed when I visit a non-existant directory
	As an author
	I want to be shown a informative error message

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: The 404 page
		Given I am on the page for non-existant directory "Pancakes"
		Then I should see text '404'
		And there should be an alert with the message "The item you were looking for cannot be found"