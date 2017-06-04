Then %r{^I should see a '(.*)' field with type '(.*)'$} do |name, field_type|
  expect(page).to have_css("input[name='#{name.downcase}'][type='#{field_type}']")
end

Then %r{^the submit button should be labelled '(.*)'$} do |label|
  expect(page).to have_css("input.btn[value='#{label}']")
end

When %r{^I enter '(.*)' in the '(.*)' field$} do |value, field|
  input = page.find("label", text: /^#{field}$/)['for']

  fill_in input, with: value
end

When %r{^I submit the form$} do
  within("form") do
    find("input[type='submit']").click
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