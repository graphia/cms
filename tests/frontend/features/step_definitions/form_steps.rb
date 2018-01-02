Then %r{^I should see a '(.*)' field with type '(.*)'$} do |name, field_type|
  expect(page).to have_css("input[name='#{name.downcase}'][type='#{field_type}']")
end

Then %r{^I should see a text area called '(.*?)'} do |name|
  expect(page).to have_css("textarea[name='#{name.downcase}']")
end

Then %r{^the submit button should be labelled '(.*)'$} do |label|
  expect(page).to have_css("input.btn[value='#{label}']")
end

# FIXME standardise on single or double quotes
When %r{^I enter "(.*)" in the "(.*)" field$} do |value, field|
  input = page.find("label", text: /^#{field}$/)['for']

  fill_in input, with: value
end

# FIXME standardise on single or double quotes
When %r{^I enter '(.*)' in the '(.*)' field$} do |value, field|
  field_matcher = /^#{field}$/
  expect(page).to have_css("label", text: field_matcher)
  input = page.find("label", text: field_matcher)['for']
  fill_in input, with: value
end

Then %r{^the "(.*?)" field should be "(.*?)"$} do |field, value|
  name = page.find("label", text: /^#{field}$/)['for']
  input = page.find("input[name='#{name}']")
  expect(input.value).to eql(value)
end

Given %r{^I fill in the form with the following data:$} do |table|
  table.rows_hash.each do |field, value|
    step "I enter '#{value}' in the '#{field}' field"
  end
end

When %r{^I submit the form$} do
  within("form") do
    find("input[type='submit']").native.send_keys(:enter)
  end
end

When %r{^I submit the "(.*?)" form$} do |form_id|
  within("form##{get_class(form_id)}") do
    find("input[type='submit']").native.send_keys(:enter)
  end
end

def get_class(desc)
  desc.gsub(" ", "-").downcase
end

When %r{^I submit the form by clicking '(.*)'$} do |label_text|
  within("form") do
    page.find("input[type='submit'][value='#{label_text}']").click
  end
end

# | Name     | Type     | Required  |
# | Name     | Text     | yes       |
# | Password | Password | yes       |
Then %r{^I should see a form with the following fields:$} do |table|

  table.hashes.each do |row|

    name = row.fetch("Name")
    data_type = row.fetch("Type").downcase
    required = (row.fetch("Required", "no").downcase == "yes")
    disabled = (row.fetch("Disabled", "no").downcase == "yes")

    within("form") do
      label = page.find("label", text: /^#{name}$/)
      input = page.find("input[name='#{label['for']}']")

      # ensure the element exists and has the right attributes
      expect(input).not_to be_nil
      expect(input['type']).to eql(data_type)
      expect(input['required']).to eql('true') if required
      expect(input['disabled']).to eql('true') if disabled
    end

  end
end

# | Name     | Type     | Required  |
# | Name     | Text     | yes       |
# | Password | Password | yes       |
Then %r{^I should see a "(.*?)" form with the following fields:$} do |form_id, table|

    table.hashes.each do |row|

      name = row.fetch("Name")
      data_type = row.fetch("Type").downcase
      required = (row.fetch("Required", "no").downcase == "yes")
      disabled = (row.fetch("Disabled", "no").downcase == "yes")

      within("form##{get_class(form_id)}") do
        label = page.find("label", text: /^#{name}$/)
        input = page.find("input[name='#{label['for']}']")

        # ensure the element exists and has the right attributes
        expect(input).not_to be_nil
        expect(input['type']).to eql(data_type)
        expect(input['required']).to eql('true') if required
        expect(input['disabled']).to eql('true') if disabled
      end

    end
  end

Then %r{^the '(.*)' field should allow values from '(\d+)' to '(\d+)' characters$} do |field, min, max|

  # this is a bit hacky, but it allows capybara to wait for the page to have loaded fully
  # by ensuring that we don't continue until there's at least one field on the page
  # that has a minlength or maxlength attr
  expect(page).to have_xpath("//*[@minlength|@maxlength]")

  name = page.find("label", text: /^#{field}$/)['for']
  input = page.find("input[name='#{name}']")
  expect(input['minlength']).to eql(min)
  expect(input['maxlength']).to eql(max)
end

And %r{^the '(.*)' field should be at least '(\d+)' characters long$} do |field, min|
  name = page.find("label", text: /^#{field}$/)['for']
  input = page.find("input[name='#{name}']")
  expect(input['minlength']).to eql(min)
end

And %r{^the '(.*)' field should be at most '(\d+)' characters long$} do |field, max|
  within("form") do
    name = page.find("label", text: /^#{field}$/)['for']
    input = page.find("input[name='#{name}']")
    expect(input['maxlength']).to eql(max)
  end
end

When %r{^I enter a '(\d+)' letter word into '(.*)'$} do |chars, field|
  val = 'a' * chars.to_i
  fill_in(field.downcase, with: val)
  expect(page.find("input[name='#{field.downcase}']").value).to eql(val)
  page.find('body').click
end

Then %r{^the "(.*?)" submit button should be enabled$} do |form_id|
  within("form##{get_class(form_id)}") do
    expect(page.find("input[type='submit']")).not_to be_disabled
  end
end

Then %r{^the "(.*?)" submit button should be disabled$} do |form_id|
  within("form##{get_class(form_id)}") do
    expect(page.find("input[type='submit']")).to be_disabled
  end
end

When %r{^I check the "(.*)" checkbox$} do |checkbox|
  check(checkbox)
end

When %r{^I uncheck the "(.*)" checkbox$} do |checkbox|
  uncheck(checkbox)
end

Then %r{^there should be a checkbox called "(.*?)"$} do |name|
  # As we're going with a label containing a checkbox for the most part
  # expect(page).to have_css("input[type='checkbox'][name='#{name.downcase}']")

  within("label", text: /^#{name}$/) do
    expect(page).to have_css("input[type='checkbox']")
  end
end