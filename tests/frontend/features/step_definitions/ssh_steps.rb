def user_with_key
  return {
    username: "rod.flanders",
    name: "Rod Flanders",
    email: "rod.flanders@springfield-elementary.k12.us",
    password: "okily-dokily!"
  }
end

Given %r{^my user account with public key exists$} do
  uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Post.new(uri, "Content-Type" => "application/json")
  req.body = user_with_key.to_json
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPCreated)
end

Given %r{^my private key is valid$} do
  pending # Write code here that turns the phrase above into concrete actions
end

When %r{^I initiate a SSH connection to the server$} do
  pending # Write code here that turns the phrase above into concrete actions
end

Then %r{^I should see the response "([^"]*)"$} do |arg1|
  pending # Write code here that turns the phrase above into concrete actions
end