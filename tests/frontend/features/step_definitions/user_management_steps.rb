Given %r{^that '(\d+)' extra users have been created$} do |number|

  # first get the token we've logged in with
  visit ("/cms/")
  token = evaluate_script("localStorage.token")

  1.upto(number.to_i).each do |i|

    uri = URI('http://127.0.0.1:9095/api/admin/users')
    req = Net::HTTP::Post.new(uri, "Content-Type" => "application/json")
    req.body = {
      username: "user.#{i}",
      name: "User Number #{i}",
      email: "user.number#{i}@somecompany.com"
    }.to_json

    req['Authorization'] = "Bearer #{token}"

    res = Net::HTTP.start(uri.hostname, uri.port) do |http|
      http.request(req)
    end

    expect(res.class.name).to eql('Net::HTTPCreated')

  end
end

When %r{^I am on the users list page$} do
  path = "/cms/settings/users"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a user list with '(\d+)' entries$} do |number|
  within(".user-list") do
    expect(page).to have_css(".user", count: number)
  end
end

Then %r{^I should see a section with my user's name as the title$} do
  within(".user-list") do
    expect(page).to have_css(".card-header", text: "Rod Flanders")
  end
end

Then %r{^the details listed should be:$} do |table|
  within(".user-list") do
    table.rows.to_h.each do |attribute, value|
      expect(page).to have_css('dt', attribute)
      expect(page).to have_css('dd', value)
    end
  end
end