Feature: Listing documents
	So I can see files
	As an author
	I want to view listings for a directory

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: Documents are visible on the documents page
		Given there are directories called "documents" and "appendices"
		And they both contain Markdown files
		When I navigate to the "documents" index page
		Then I should see a list containing the contents of the "documents" directory

	Scenario: Documents are visible on the documents page
		Given I am on the "documents" index page
		When I click the "Appendices" navigation link
		Then I should be on the "appendices" index page
		And I should see a list containing the contents of the "appendices" directory

	Scenario: The page title should match the directory name
		Given there are directories called "documents" and "appendices"
		Then each directory index page should have the correct title:
			| Directory  | Title      |
			| documents  | Documents  |
			| appendices | Appendices |

	Scenario: Breadcrumbs
		Given I am on the "documents" index page
		Then I should see the following breadcrumbs:
			| Text                | Reference |
			| Dashboard           | /cms      |
			| Important Documents | None      |