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