Given %r{^a user account has been created$} do
  visit('/') # this lets us wait for the web server to start up
  user = {
    username: "rod.flanders",
    name: "Rod Flanders",
    email: "rod.flanders@springfield-elementary.k12.us",
    password: "okily-dokily!"
  }
  uri = URI('http://127.0.0.1:9095/auth/create_initial_user')
  req = Net::HTTP::Post.new(uri, "Content-Type" => "application/json")
  req.body = user.to_json
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
end

Given %r{I am on the login screen} do
  path = "/cms/login"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a '(.*)' field with type '(.*)'$} do |name, field_type|
  expect(page).to have_css("input[name='#{name.downcase}'][type='#{field_type}']")
end

Then %r{^the submit button should be labelled '(.*)'$} do |label|
  expect(page).to have_css("input.btn[value='#{label}']")
end

When %r{^I enter invalid credentials$} do
  fill_in :username, with: "artie.ziff"
  fill_in :password, with: "Zifc0rP"
end

When %r{^I submit the form$} do
  within("form") do
    find("input[type='submit']").click
  end
end

Then %r{^I should still be on the login screen$} do
  path = "/cms/login"
  expect(page.current_path).to eql(path)
end

Then %r{^there should be a 'red' alert with the message 'Invalid'$} do
end