Feature: Home page contents

	So I can see an overview of the contents
	As a user
	I want the homepage include a statistics panel

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: The panel should be visible on the homepage
		Given I am on the homepage
		Then I should see a "Statistics" section

	Scenario: The panel should contain correct user counts
		Given there is one user
		When I am on the homepage
		Then the statistics panel's "user" count should equal "1"

	Scenario: The panel should contain correct commit counts
		Given there is one user
		And I have made changes to an existing file
		When I am on the homepage
		Then the statistics panel's "commit" count should equal "3"

	Scenario: The panel should contain the correct file counts
		Given I am on the homepage
		Then I should see "8" files of type "Documents"
		And I should see "4" files of type "Images"
		And I should see "2" files of type "Structured Data"
		And I should see "1" files of type "Tabular Data"
		And I should see "2" files of type "Other"