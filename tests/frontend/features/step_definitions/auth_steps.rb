Given %r{^I am not logged in$} do
  # do nothing
end

When %r{^I navigate directly to the homepage$} do
  visit("/cms/home")
end

Given %r{^I navigate directly to a protected page$} do
  @original_path = "/cms/documents"
  visit(@original_path)
end

Then %r{^I should be redirected to the login page$} do
  expect(page.current_path).to eql("/cms/login")
end

When %r{^I am prompted for my credentials$} do
  step "I should be redirected to the login page"
end

When %r{^I try to manually make an unauthenticated HTTP request to the API$} do
  uri = URI('http://127.0.0.1:9095/api/directories/documents/documents')
  req = Net::HTTP::Get.new(uri, "Content-Type" => "application/json")
  @res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
end

When %r{^I try to manually make an authenticated HTTP request to the API$} do
  token = evaluate_script("localStorage.token")
  uri = URI('http://127.0.0.1:9095/api/directories/documents/documents')
  req = Net::HTTP::Get.new(
    uri,
    "Content-Type" => "application/json",
    "Authorization" => "Bearer #{token}"
  )
  @res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
end

Then %r{^I should receive a 'HTTP Not Authorized' status$} do
  expect(@res.code_type).to eql(Net::HTTPUnauthorized)
end

Then %r{^I should receive a 'HTTP OK' status$} do
  expect(@res.code_type).to eql(Net::HTTPOK)
end

When %r{^I provide them and log in$} do
  fill_in :username, with: user[:username]
  fill_in :password, with: user[:password]
  step "I submit the form by clicking 'Log in'"
end

Then %r{^I should be authenticated and redirected to my original destination$} do
  expect(page).to have_content("You have logged in successfully")
  expect(page.current_path).to eql(@original_path)
end