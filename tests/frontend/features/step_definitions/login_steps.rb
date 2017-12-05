def user
  return {
    username: "rod.flanders",
    name: "Rod Flanders",
    email: "rod.flanders@springfield-elementary.k12.us",
    password: "okily-dokily!"
  }
end

Given %r{^(?:a user account has been created|my user account exists)$} do
  uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Post.new(uri, "Content-Type" => "application/json")
  req.body = user.to_json
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPCreated)
end

Given %r{I am on the login screen} do
  path = "/cms/login"
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I enter invalid credentials$} do
  fill_in :username, with: "artie.ziff"
  fill_in :password, with: "Zifc0rP"
end

When %r{^I enter valid credentials$} do
  fill_in :username, with: user[:username]
  fill_in :password, with: user[:password]
end

Then %r{^I should be redirected to the CMS's landing page$} do
  expect(page.current_path).to eql("/cms")
end

Then %r{^I should still be on the login screen$} do
  path = "/cms/login"
  expect(page.current_path).to eql(path)
end

Given %r{^I have logged in$} do
  steps %{
		Given I am on the login screen
		When I enter valid credentials
		And I submit the form
		Then I should see a message containing 'You have logged in successfully'
		And I should be redirected to the CMS's landing page
  }
end

Then %r{^I should have a JWT saved in localstorage$} do
  token = evaluate_script("localStorage.token")
  expect(token.split(".").size).to eql(3) # it at least looks like a JWT!
end

Given %r{^I am logged in$} do
  step %{I have logged in}
  @get_token = "localStorage.getItem('token');"
  expect(page.evaluate_script(@get_token)).not_to be_nil
end

When %r{^I select 'Logout' from the settings menu$} do
  within("nav.navbar") do
    page.find("#user-menu").click
    page.find(".dropdown-item.logout").click
  end
end

Then %r{^I should be logged out$} do
  expect(page.evaluate_script(@get_token)).to be_nil
end

Then %r{^I should be on the login screen$} do
  expect(page.current_path).to eql("/cms/login")
end