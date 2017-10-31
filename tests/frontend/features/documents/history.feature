Feature: Commits
	So I can review a file's history
	As an author
	I want to be able to view every revision made to it

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Viewing a document's history
		Given there is a document with '3' revisions
		When I am on the document's history page
		Then I should see each revision listed

	Scenario: Viewing a revision's commit
		Given there is a document with some revisions
		And I am on the document's history page
		When I click the 'View entire commit' button for a particular revision
		Then I should be on that commit's page

	Scenario: Viewing a revision in-line
		Given there is a document with some revisions
		And I am on the document's history page
		When I click the 'Show changes' button for a particular revision
		Then that revision's diff should be visible beneath the revision entry
		And the diff should have correctly marked insertions and deletions

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