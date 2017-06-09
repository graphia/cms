Given %r{^I am not logged in$} do
  # do nothing
end

When %r{^I navigate directly to the homepage$} do
  path = "/cms/home"
  visit(path)
end

Then %r{^I should be redirected to the login page$} do
  expect(page.current_path).to eql("/cms/login")
end

When %r{^I try to manually make an unauthenticated HTTP request to the API$} do
  uri = URI('http://127.0.0.1:9095/api/directories/documents/files')
  req = Net::HTTP::Get.new(uri, "Content-Type" => "application/json")
  @res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
end

When %r{^I try to manually make an authenticated HTTP request to the API$} do
  token = evaluate_script("localStorage.token")
  uri = URI('http://127.0.0.1:9095/api/directories/documents/files')
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