Feature: Home page contents

	So I can see an overview of the CMS
	As a user
	I want the homepage to provide me with useful information

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page heading
		Given I am on the homepage
		Then the main heading should be "Dashboard"

	Scenario: Page title is dynamic
		Given I am on the homepage
		And the page's title should be "Graphia CMS"

	Scenario: Primary navigation bar should contain directory links
		Given the following directories exist in the repository
			| appendices |
			| documents  |
		When I am on the homepage
		Then the navigation bar should contain the following links:
			| Appendices           |
			| Important Documents  |

	Scenario: Home page sections
		Given I am on the homepage
		Then I should see a summary of recent changes
		And I should see a statistics section

	Scenario: Recent commits
		Given there have been some recent commits
		And I am on the homepage
		Then the recent changes summary should contain a list of commits

	Scenario: Listing directories
		Given the following directories exist in the repository
			| appendices |
			| documents  |
		When I am on the homepage
		Then I should see a section for each directory

	Scenario: Files within a repository
		Given the documents directory contains the following files:
			| document_1 |
			| document_2 |
			| document_3 |
		When I am on the homepage
		Then I should see all three documents listed
		And there should be a 'New document' button

	Scenario: When a directory is empty
		Given the 'empty' directory contains no files
		When I am on the homepage
		Then I see a 'no files' alert in the empty section
		And there should be a 'New document' button

	Scenario: When a directory has metadata
		Given the 'documents' directory has title and description metadata
		When I am on the homepage
		Then I should see the custom description
		And I should see the custom title

	Scenario: Breadcrumbs
		Given I am on the homepage
		Then I should only see the inactive 'Dashboard' breadcrumb

	Scenario: Primary navigation highlighting
		Given I am on the homepage
		Then the primary navigation link "Home" should be active
		When I am on the "documetns" index page
		Then the primary navigation link "Home" should not be active