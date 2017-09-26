Feature: Commits
	So I can review a file's history
	As an author
	I want to be able to view every revision made to it

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in

	Scenario: Breadcrumbs without metadata
		Given I am on the appendix history page for "appendix_1.md"
		Then I should see the following breadcrumbs:
			| Text                | Reference                     |
			| Dashboard           | /cms                          |
			| appendices          | /cms/appendices               |
			| appendix_1.md       | /cms/appendices/appendix_1.md |
			| History             | None                          |

	Scenario: Breadcrumbs with metadata
		Given I am on the document history page for "document_1.md"
		Then I should see the following breadcrumbs:
			| Text                | Reference                    |
			| Dashboard           | /cms                         |
			| Important Documents | /cms/documents               |
			| document 1          | /cms/documents/document_1.md |
			| History             | None                         |