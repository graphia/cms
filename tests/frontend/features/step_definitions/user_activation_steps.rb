Given %r{^the administrator has logged in$} do
  step "I have logged in"
end

Given %r{^there is a user requiring activation$} do
  steps %{
		And a user account has been created
    And the administrator has logged in
		And the administrator has created a new user account
		And the administrator has logged out
  }
end

# as an admin, create an account, grab the token and use it
# to get the new account's confirmation key
Given %r{^the administrator has created a new user account$} do
  steps %{
    Given I am on the new user page
    When I fill in the form with the following data:
      | Username         | hhermann       |
      | Name             | Herman Hermann |
      | Email address    | hello@hma.com  |
    And I submit the form
    Then I should see a message containing 'Herman Hermann will receive an email with instructions on how to log in'
    And I should be on the users list page
  }

  token = evaluate_script("localStorage.token")

  # now get the confirmation key
  uri = URI('http://127.0.0.1:9095/api/admin/users/hhermann')
  req = Net::HTTP::Get.new(uri, "Content-Type" => "application/json")
  req['Authorization'] = "Bearer #{token}"
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPOK)

  @user_details = JSON.parse(res.body)
  @confirmation_key = @user_details['confirmation_key']

  expect(@confirmation_key.length).to eql(32)

end

Then %r{^the password fields should be required and at least (\d+) characters long$} do |arg1|
  up = page.find("input[name='user-password']")
  expect(up['required']).to be_truthy
  expect(up['minlength']).to eql('6')

  pc = page.find("input[name='user-confirm']")
  expect(pc['required']).to be_truthy
end

Given %r{^the administrator has logged out$} do
  steps %{
    And I select 'Logout' from the settings menu
  }
end

Given %r{^I am on the activate user screen$} do
  path = "/cms/activate/#{@confirmation_key}"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a personal welcome$} do
  expect(page).to have_css("h1", text: "Welcome #{@user_details['name']}")
end

Then %r{^there should be a password reset form$} do
  expect(page).to have_css("input[name='user-password'][type='password']")
  expect(page).to have_css("input[name='user-confirm'][type='password']")
end

Given %r{^I fill in and confirm the password$} do
  @password = "abcd1234"

  steps %{
    When I enter "#{@password}" in the "Password" field
    And I enter "#{@password}" in the "Confirm Password" field
  }
end


Given %r{^I have activated my account$} do
  steps %{
    Given I am on the activate user screen
    And I fill in and confirm the password
    When I submit the form
    Then I should see a message containing 'Your account has been activated'
    And I should be on the login screen
  }
end

When %r{^I enter my username and password$} do
  fill_in "username", with: @user_details['username']
  fill_in "password", with: @password
end

Then %r{^I should be logged in$} do
  step "I should see a message containing 'You have logged in successfully'"
  get_token = "localStorage.getItem('token');"
  expect(page.evaluate_script(get_token)).not_to be_nil
end
