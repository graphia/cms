Feature: Home page contents

	So I can see an overview of the CMS
	As a user
	I want the homepage to provide me with useful information

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: Page heading
		Given I am on the homepage
		Then the main heading should be "Dashboard"

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
			| document_1.md |
			| document_2.md |
			| document_3.md |
		When I am on the homepage
		Then I should see all three documents listed
		And there should be a 'new file' button

	Scenario: When a directory is empty
		Given the 'empty' directory contains no files
		When I am on the homepage
		Then I see a 'no files' alert in the operating procedures section
		And there should be a 'new file' button

	Scenario: When a directory has metadata
		Given the 'documents' directory has title and description metadata
		When I am on the homepage
		Then I should see the custom description
		And I should see the custom title

	Scenario: Breadcrumbs
		Given I am on the homepage
		Then I should only see the inactive 'Dashboard' breadcrumb