Then %r{^I should see a '(.*)' field with type '(.*)'$} do |name, field_type|
  expect(page).to have_css("input[name='#{name.downcase}'][type='#{field_type}']")
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
  input = page.find("label", text: /^#{field}$/)['for']

  fill_in input, with: value
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
    required =
      case row.fetch("Required").downcase
      when "yes"
        true
      else
        false
      end

    within("form") do
      label = page.find("label", text: /^#{name}$/)
      input = page.find("input[name='#{label['for']}']")

      # ensure the element exists and has the right attributes
      expect(input).not_to be_nil
      expect(input['type']).to eql(data_type)
      expect(input['required']).to eql('true') if required
    end


  end
end

Then %r{^the '(.*)' field should allow values from '(\d+)' to '(\d+)' characters$} do |field, min, max|
  within("form") do
    name = page.find("label", text: /^#{field}$/)['for']
    input = page.find("input[name='#{name}']")
    expect(input['minlength']).to eql(min)
    expect(input['maxlength']).to eql(max)
  end
end

And %r{^the '(.*)' field should be at least '(\d+)' characters long$} do |field, min|
  within("form") do
    name = page.find("label", text: /^#{field}$/)['for']
    input = page.find("input[name='#{name}']")
    expect(input['minlength']).to eql(min)
  end
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