Then %r{^I should see a button with text '(.*?)'$} do |caption|
  expect(page).to have_css("button.btn-primary", text: caption)
end

Then %r{^I should see text '(.*?)'$} do |text|
  expect(page).to have_content(text)
end

Then %r{^there should be an alert with the message "(.*)"$} do |message|
  expect(page).to have_css(".alert", text: message)
end

Then %r{^the main heading should be "(.*?)"$} do |heading|
  expect(page.first("h1")).to have_content(heading)
end
