Given %r{^there are no users$} do
end

When %r{^I navigate to the login page$} do
  path = '/cms/login'
  visit(path)
end

When %r{^I navigate to the setup page$} do
  path = '/cms/setup/initial_user'
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should be redirected to the initial setup page$} do
  expect(page).to have_css("label", text: "Username")
  expect(page.current_path).to eql('/cms/setup/initial_user')
end

Given %r{^I am on the initial setup page$} do
  steps %{
      Given there are no users
      And I navigate to the setup page
  }
end

Then %r{^the '(.*)' field should be invalid$} do |field|
  selector = %Q{$("form input[name='#{field.downcase}']").get(0).checkValidity() }
  valid = evaluate_script(selector)
  expect(valid).to be false
end

Then %r{^the '(.*)' field should be valid$} do |field|
  selector = %Q{$("form input[name='#{field.downcase}']").get(0).checkValidity() }
  valid = evaluate_script(selector)
  expect(valid).to be true
end

Then %r{^the 'Confirm Password' field should be marked invalid$} do
  within("form .confirm-password-group") do
    expect(page).to have_css("input.is-invalid")
  end
end

Then %r{^there should be a warning containing '(.*)'$} do |warning|
  within("form .confirm-password-group") do
    expect(page).to have_content(warning)
    expect(page).to have_css(".password-match-feedback.form-control-feedback-message")
  end
end

When %r{^I enter matching passwords in the '(.*)' and '(.*)' fields$} do |first, second|
  [first, second].each do |name|
    step "I enter 'password' in the '#{name}' field"
  end
end

Then %r{^no password\-related warnings should be visible$} do
  within("form .confirm-password-group") do
    expect(page).not_to have_css("input.is-invalid")
    expect(page).not_to have_css(".password-match-feedback.form-control-feedback-message")
  end
end

Then %r{^the new user should have been saved to the database$} do
  uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Get.new(uri, "Content-Type" => "application/json")
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
  json = JSON.parse(res.body)
  expect(json['enabled']).to eql(false)
end

Then %r{^I should see a message containing '(.*)'$} do |message|
  within(".row.messages") do
    expect(page).to have_css("div.alert.alert-success", text: /#{message}/)
  end
end
