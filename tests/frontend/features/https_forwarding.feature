Feature: HTTPS Forwarding

	So I can ensure the site is displayed as intended
	As a user
	I want my connection to the server to be encrypted

	Background:
		Given a repository has been initialised
		And the CMS is running with the "https_enabled" config

	Scenario Outline: Ensure forwarding works
		Given I try to access "<endpoint>" via HTTP
		Then I should receive a redirect to the same "<endpoint>" but served via HTTPS

		Examples:
			| endpoint                   |
			| cms/appendices             |
			| cms/appendices/appendix_1  |
			| api/settings/ssh           |
			| api/directories            |
			| auth/login                 |
			| setup/create_initial_user  |