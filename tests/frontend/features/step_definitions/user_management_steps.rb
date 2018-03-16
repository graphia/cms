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

Given %r{^I am on the new user page$} do
  path = "/cms/settings/users/new"
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
      expect(page).to have_css('dt', text: attribute)
      expect(page).to have_css('dd', text: value)
    end
  end
end

Then %r{^there should be a user list$} do
  expect(page).to have_css(".user-list")
end

When %r{^I click the "(.*?)" button for my user$} do |caption|
  within("#user-1") do
    expect(page).to have_css(".card-header", text: "Rod Flanders")
    click_link "Edit"
  end
end

Then %r{^I should be on the edit page for my user$} do
  expect(page.current_path).to eql("/cms/settings/users/rod.flanders/edit")
end

When %r{^I click the "(.*?)" link$} do |caption|
  click_link caption
end

Then %r{^I should be on the new user page$} do
  expect(page.current_path).to eql("/cms/settings/users/new")
end

Then %r{^I should be on the users list page$} do
  expect(page.current_path).to eql("/cms/settings/users")
end


Then %r{^I see my newly\-created user when redirected to the user list$} do
  expect(page.current_path).to eql("/cms/settings/users")
  expect(page).to have_content("Herman Hermann")
  expect(page).to have_content("hhermann")
  expect(page).to have_content("hello@hma.com")

end


  Given %r{^there is a regular user and an administrator$} do
  # admin is already created, create a regular user

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
  end

Then %r{^the 'administrator' should have an 'Admin' label$} do
  within("#user-1") do
    expect(page).to have_css("span.admin", text: 'Admin')
  end
end

Then %r{^the 'regular user' should have no 'Admin' label$} do
  within("#user-2") do
    expect(page).not_to have_css("span.admin")
  end
end

Given %r{^I am on the edit user page$} do
  path = "/cms/settings/users/rod.flanders/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I re\-enter the details of an existing user$} do
  steps %{
    When I fill in the form with the following data:
      | Username         | rod.flanders                                |
      | Name             | Rod Flanders                                |
      | Email address    | rod.flanders@springfield-elementary.k12.us  |
  }
end

When %r{^I change the name and email address to match the regular user$} do
  steps %{
    When I fill in the form with the following data:
      | Name             | Herman Hermann |
      | Email address    | hello@hma.com  |
  }
end

Then %r{^I should see an error message stating that the record is not unique$} do
  msg = "This record cannot be saved because either the username or email address are already in use"
  expect(page).to have_css(".alert.alert-danger", text: msg)
end

Then %r{^I see my newly\-updated user when redirected to the user list$} do
  within("#user-1") do
    expect(page).to have_content("Todd Flanders")
    expect(page).to have_content("tf@floody.com")
  end
end

Given %r{^I am on the edit user page for the regular user$} do
  path = "/cms/settings/users/hhermann/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on my own edit user page$} do
  path = "/cms/settings/users/rod.flanders/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^the regular user should not be present$} do
  # check the admin is present too so we know
  # we're using the right selector
  expect(page).to have_css(".card", text: "Rod Flanders")
  expect(page).not_to have_css(".card", text: "Herman Hermann")
end

Given %r{^I am on the edit page for non\-existing user '(.*?)'$} do |username|
  path = "/cms/settings/users/#{username}/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end