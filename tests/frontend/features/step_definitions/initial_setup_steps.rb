Given %r{^there are no users$} do
end

When %r{^I navigate to the login page$} do
  path = '/cms/login'
  visit(path)
end

When %r{^I navigate to the setup page$} do
  path = '/cms/setup'
  visit(path)
  expect(page.current_path).to eql('/cms/setup')
end

Then %r{^I should be redirected to the initial setup page$} do
  expect(page.current_path).to eql('/cms/setup')
end

Given %r{^I am on the initial setup page$} do
  steps %{
      Given there are no users
      And I navigate to the setup page
  }
end

When %r{^I enter a '(\d+)' letter word into '(.*)'$} do |chars, field|

  val = 'a' * chars.to_i
  fill_in(field.downcase, with: val)
  expect(page.find("input[name='#{field.downcase}']").value).to eql(val)
  page.find('body').click
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
  within("form") do
    expect(page).to have_css(".confirm-password-group.has-danger")
  end
end

Then %r{^there should be a warning containing '(.*)'$} do |warning|
  within("form .confirm-password-group") do
    expect(page).to have_content(warning)
    expect(page).to have_css(".passwords-do-not-match-message")
  end
end

When %r{^I enter matching passwords in the '(.*)' and '(.*)' fields$} do |first, second|
  [first, second].each do |name|
    step "I enter 'password' in the '#{name}' field"
  end
end

Then %r{^no password\-related warnings should be visible$} do
  within("form") do
    expect(page).not_to have_css(".confirm-password-group.has-danger")
    expect(page).not_to have_css(".passwords-do-not-match-message")
  end
end