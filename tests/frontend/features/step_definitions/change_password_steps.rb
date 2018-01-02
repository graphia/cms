Then %r{^I should see a section called "([^"]*)"$} do |text|
  within(".password-update") do
    expect(page).to have_css("h2", text: text)
  end
end

# FIXME this clashes with an outdated single-quoted version
Then %r{^the "Confirm Password" field should be marked invalid$} do
  expect(page).to have_css("input.is-invalid[name='confirmPassword']")
end

Then %r{^the "New Password" and "Confirm Password" fields should be marked valid$} do
  expect(page).to have_css("input.is-valid[name='newPassword']")
  expect(page).to have_css("input.is-valid[name='confirmPassword']")
end

Then %r{^I should see an error containing "([^"]*)"$} do |message|
  expect(page).to have_css(".alert.alert-danger", text: message)
end