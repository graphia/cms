Given %r{^a user account has been created$} do
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

Then %r{^I should still be on the login screen$} do
  path = "/cms/login"
  expect(page.current_path).to eql(path)
end

Then %r{^there should be a 'red' alert with the message 'Invalid'$} do
end