Then %r{^I should see a button with text '(.*?)'$} do |caption|
  expect(page).to have_css("button.btn-primary", text: caption)
end